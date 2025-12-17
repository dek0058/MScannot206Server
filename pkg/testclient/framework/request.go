package framework

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func WebRequest[ReqT any, ResT any](c Client) webRequest[ReqT, ResT] {
	return webRequest[ReqT, ResT]{
		client:  c,
		headers: make(map[string]string, 8),
	}
}

type webRequest[ReqT any, ResT any] struct {
	client Client

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

	baseURL := strings.TrimRight(r.client.GetUrl(), "/")
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

	resp, err := r.client.Do(req)
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

	resp, err := r.client.Do(req)
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
