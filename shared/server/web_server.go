package server

import (
	"MScannot206/shared/config"
	"MScannot206/shared/service"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"
)

func NewWebServer(
	ctx context.Context,
	cfg *config.WebServerConfig,
) (*WebServer, error) {
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
		services: make([]service.GenericService, 0),
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

	services []service.GenericService
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

func (s *WebServer) Init() error {
	var errs error
	var wg sync.WaitGroup
	taskCh := make(chan error, 1)
	tasks := []func(context.Context, chan error){
		s.connectMongoTask(&wg, taskCh),
	}

	for _, task := range tasks {
		wg.Add(1)
		go task(s.ctx, taskCh)
	}

	allDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(allDone)
	}()

Loop:
	for {
		select {
		case taskErr := <-taskCh:
			if taskErr != nil {
				errs = errors.Join(errs, taskErr)
				log.Printf("초기화 에러 발생: %v", taskErr)
			}

		case <-allDone:
			log.Printf("초기화 작업 완료")
			break Loop

		case <-s.ctx.Done():
			log.Printf("서버 Context 취소, 초기화 중단")
			return s.ctx.Err()
		}

	}

	if errs != nil {
		return errs
	}

	errs = nil
	for _, svc := range s.services {
		if err := svc.Init(); err != nil {
			errs = errors.Join(errs, err)
			log.Println(err)
		}
	}

	if errs != nil {
		return errs
	}

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
			log.Println(err)
		}
	}

	return s.server.ListenAndServe()
}

func (s *WebServer) Quit() error {
	for _, svc := range s.services {
		if err := svc.Stop(); err != nil {
			log.Println(err)
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

func (s WebServer) GetServices() []service.GenericService {
	return s.services
}

func (s *WebServer) AddService(svc service.GenericService) error {
	if svc == nil {
		return errors.New("service is null")
	}

	s.services = append(s.services, svc)
	return nil
}

func (s *WebServer) connectMongoTask(
	wg *sync.WaitGroup, taskCh chan error,
) func(context.Context, chan error) {
	return func(ctx context.Context, errCh chan error) {
		defer wg.Done()

		var err error

		connectCtx, connectCancel := context.WithTimeout(ctx, 10*time.Second)
		defer connectCancel()

		opts := mongo_options.Client().ApplyURI(s.cfg.MongoUri)
		s.mongoClient, err = mongo.Connect(connectCtx, opts)

		if err != nil {
			log.Printf("MongoDB 연결 실패[uri:%v][err:%v]", s.cfg.MongoUri, err)
			taskCh <- err
			return
		}

		log.Printf("MongoDB 연결 완료[uri:%v]", s.cfg.MongoUri)
	}
}
