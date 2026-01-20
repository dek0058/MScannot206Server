package server

import (
	"MScannot206/shared/config"
	"MScannot206/shared/def"
	"MScannot206/shared/service"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrServeMuxIsNil = errors.New("serve mux is null")

func NewWebServer(
	ctx context.Context,
	cfg *config.WebServerConfig,
	mongoClient *mongo.Client,
) (*WebServer, error) {
	if mongoClient == nil {
		return nil, errors.New("mongo client is null")
	}

	ctxWithCancel, cancel := context.WithCancel(ctx)

	webServerCfg := cfg
	if webServerCfg == nil {
		webServerCfg = &config.WebServerConfig{
			Port: 8080,

			MongoUri: "mongodb://localhost:27017/",
		}
	}

	server := &WebServer{
		ctx:        ctxWithCancel,
		cancelFunc: cancel,

		cfg: webServerCfg,

		router:   http.NewServeMux(),
		services: make([]service.Service, 0),

		mongoClient: mongoClient,
	}

	return server, nil
}

type WebServer struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	// Config
	cfg *config.WebServerConfig

	// Core
	router *http.ServeMux
	server *http.Server

	// DB
	mongoClient *mongo.Client

	services []service.Service
}

func (s WebServer) GetContext() context.Context {
	return s.ctx
}

func (s WebServer) GetRouter() *http.ServeMux {
	return s.router
}

func (s WebServer) GetMongoClient() *mongo.Client {
	return s.mongoClient
}

func (s WebServer) GetLocale() def.Locale {
	return def.Locale(s.cfg.Locale)
}

func (s *WebServer) Init() error {
	addr := fmt.Sprintf(":%v", s.cfg.Port)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return nil
}

func (s *WebServer) Start() error {
	for _, svc := range s.services {
		if err := svc.Start(); err != nil {
			log.Err(err).Msg("서버 시작 중 에러가 발생하였습니다.")
		}
	}

	return s.server.ListenAndServe()
}

func (s *WebServer) Quit() error {
	for _, svc := range s.services {
		if err := svc.Stop(); err != nil {
			log.Err(err).Msg("서버 종료 중 에러가 발생하였습니다.")
		}
	}

	if err := s.server.Shutdown(s.ctx); err != nil {
		return err
	}

	if err := s.mongoClient.Disconnect(s.ctx); err != nil {
		return err
	}

	if s.cancelFunc != nil {
		s.cancelFunc()
	}

	return nil
}

func (s WebServer) GetServices() []service.Service {
	return s.services
}

func (s *WebServer) AddService(svc service.Service) error {
	if svc == nil {
		return errors.New("service is null")
	}

	s.services = append(s.services, svc)
	return nil
}
