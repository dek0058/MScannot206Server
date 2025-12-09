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
	return &WebServer{
		ctx: ctx,
		cfg: cfg,

		router: http.NewServeMux(),

		services: make([]service.GenericService, 0),
	}, nil
}

type WebServer struct {
	ctx context.Context

	// Config
	cfg *config.WebServerConfig

	// DB
	client *mongo.Client

	// Core
	router *http.ServeMux
	server *http.Server

	services []service.GenericService
}

func (s WebServer) GetContext() context.Context {
	return s.ctx
}

func (s WebServer) GetRouter() *http.ServeMux {
	return s.router
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

	// TODO: 고루틴 이용 할 수 있도록 해야 함.
	return s.server.ListenAndServe()
}

func (s *WebServer) Shutdown() error {
	for _, svc := range s.services {
		if err := svc.Stop(); err != nil {
			log.Println(err)
		}
	}

	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
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
		s.client, err = mongo.Connect(connectCtx, opts)

		if err != nil {
			log.Printf("MongoDB 연결 실패[uri:%v][err:%v]", s.cfg.MongoUri, err)
			taskCh <- err
			return
		}

		log.Printf("MongoDB 연결 완료[uri:%v]", s.cfg.MongoUri)
	}
}
