package client

import (
	"MScannot206/shared/config"
	"MScannot206/shared/service"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const HTTP_TIMEOUT = 30 * time.Second

func NewWebClient(ctx context.Context, cfg *config.WebClientConfig) (*WebClient, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}

	ctxWithCancel, cancel := context.WithCancel(ctx)

	client := http.Client{
		Timeout: HTTP_TIMEOUT,
	}

	webClientCfg := cfg
	if webClientCfg == nil {
		webClientCfg = &config.WebClientConfig{
			Url:  "http://localhost",
			Port: 8080,
		}
	}

	self := &WebClient{
		ctx:        ctxWithCancel,
		cancelFunc: cancel,

		client: &client,

		cfg: webClientCfg,

		services: make([]service.GenericService, 0),
	}

	return self, nil
}

type WebClient struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	// Config
	cfg *config.WebClientConfig
	url string

	// Core
	client *http.Client

	services []service.GenericService
}

func (c *WebClient) GetContext() context.Context {
	return c.ctx
}

func (c WebClient) GetMongoClient() *mongo.Client {
	return nil
}

func (c *WebClient) Init() error {
	var errs error

	c.url = c.cfg.Url + ":" + fmt.Sprintf("%v", c.cfg.Port)
	if c.url == "" {
		return errors.New("웹 클라이언트 URL이 비어있습니다")
	}

	for _, svc := range c.services {
		if err := svc.Init(); err != nil {
			errs = errors.Join(errs, err)
			log.Println(err)
		}
	}

	if errs != nil {
		return errs
	}

	return nil
}

func (c *WebClient) Start() error {
	for _, svc := range c.services {
		if err := svc.Start(); err != nil {
			return err
		}
	}

	<-c.ctx.Done()
	return nil
}

func (c *WebClient) Quit() error {
	for _, svc := range c.services {
		if err := svc.Stop(); err != nil {
			return err
		}
	}

	if c.cancelFunc != nil {
		c.cancelFunc()
	}

	return nil
}

func (c WebClient) GetServices() []service.GenericService {
	return c.services
}

func (s *WebClient) AddService(svc service.GenericService) error {
	if svc == nil {
		return errors.New("service is null")
	}

	s.services = append(s.services, svc)
	return nil
}

func WebRequest[ReqT any, ResT any](c *WebClient) webRequest[ReqT, ResT] {
	return webRequest[ReqT, ResT]{
		client:  c,
		headers: make(map[string]string, 8),
	}
}

type webRequest[ReqT any, ResT any] struct {
	client *WebClient

	endpoint    string
	headers     map[string]string
	pathParams  map[string]string
	queryParams map[string]string
	body        *ReqT
}

func (r webRequest[ReqT, ResT]) Endpoint(endpoint string) webRequest[ReqT, ResT] {
	r.endpoint = endpoint
	return r
}

func (r webRequest[ReqT, ResT]) Header(key, value string) webRequest[ReqT, ResT] {
	r.headers[key] = value
	return r
}

func (r webRequest[ReqT, ResT]) Body(body *ReqT) webRequest[ReqT, ResT] {
	r.body = body
	return r
}

func (r webRequest[ReqT, ResT]) Path(key, value string) webRequest[ReqT, ResT] {
	r.pathParams[key] = value
	return r
}

func (r webRequest[ReqT, ResT]) Query(key, value string) webRequest[ReqT, ResT] {
	r.queryParams[key] = value
	return r
}

func (r webRequest[ReqT, ResT]) buildURL() (string, error) {
	finalPath := r.endpoint
	for key, value := range r.pathParams {
		placeholder := "{" + key + "}"
		finalPath = strings.ReplaceAll(finalPath, placeholder, value)
	}

	baseURL := strings.TrimRight(r.client.url, "/")
	if !strings.HasPrefix(finalPath, "/") && finalPath != "" {
		finalPath = "/" + finalPath
	}
	fullURL := baseURL + finalPath

	if len(r.queryParams) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return "", err
		}
		q := u.Query()
		for key, value := range r.queryParams {
			q.Add(key, value)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	return fullURL, nil
}

func (r webRequest[ReqT, ResT]) Post() (*ResT, error) {
	if r.client == nil {
		return nil, errors.New("client is nil")
	}

	targetURL, err := r.buildURL()
	if err != nil {
		return nil, fmt.Errorf("failed to build url: %w", err)
	}

	jsonData, err := json.Marshal(r.body)
	if err != nil {
		return nil, err
	}

	requestBody := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest("POST", targetURL, requestBody)
	if err != nil {
		return nil, err
	}

	r.headers["Content-Type"] = "application/json"
	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	resp, err := r.client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("unexpected Content-Type: %s", contentType)
	}

	// TODO: 추가적인 헤더 검증

	var result ResT
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r webRequest[ReqT, ResT]) Get() (*ResT, error) {
	if r.client == nil {
		return nil, errors.New("client is nil")
	}

	targetURL, err := r.buildURL()
	if err != nil {
		return nil, fmt.Errorf("failed to build url: %w", err)
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	resp, err := r.client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("unexpected Content-Type: %s", contentType)
	}

	// TODO: 추가적인 헤더 검증

	var result ResT
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
