package main

import (
	"MScannot206/pkg/auth"
	"MScannot206/pkg/auth/session"
	"MScannot206/pkg/login"
	"MScannot206/pkg/serverinfo"
	"MScannot206/pkg/user"
	"MScannot206/shared/config"
	"MScannot206/shared/server"
	"MScannot206/shared/service"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func loadDefaultConfig(executablePath string, defaultFileName string, cfg interface{}) {
	defaultPath := filepath.Join(filepath.Dir(executablePath), defaultFileName)
	if _, err := os.Stat(defaultPath); err == nil {
		if err := config.LoadYamlConfig(defaultPath, cfg); err != nil {
			log.Warn().Err(err).Msgf("기본 설정 파일(%s) 로드 오류", defaultFileName)
		}
	} else {
		log.Warn().Msgf("기본 설정 파일(%s) 을(를) 찾을 수 없습니다", defaultFileName)
	}
}

func setupConfig(
	logCfg *config.LogConfig,
	serverCfg *config.WebServerConfig) error {

	var errs error
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	var logCfgPath = flag.String("logconfig", "", "로그 설정 파일 경로 지정")
	var serverCfgPath = flag.String("serverconfig", "", "서버 설정 파일 경로 지정")
	flag.Parse()

	// 로그 설정 로드
	if *logCfgPath != "" {
		if err := config.LoadYamlConfig(*logCfgPath, logCfg); err != nil {
			loadDefaultConfig(ex, "server_log_config.yaml", logCfg)
		}
	} else {
		loadDefaultConfig(ex, "server_log_config.yaml", logCfg)
	}

	// 서버 설정 로드
	if *serverCfgPath != "" {
		if err := config.LoadYamlConfig(*serverCfgPath, serverCfg); err != nil {
			loadDefaultConfig(ex, "server_config.yaml", serverCfg)
		}
	} else {
		loadDefaultConfig(ex, "server_config.yaml", serverCfg)
	}

	return errs
}

func setupServices(server *server.WebServer, cfg *config.WebServerConfig) error {
	var errs error

	// 서버정보 서비스
	serverInfo_service, err := serverinfo.NewServerInfoService(server, cfg.ServerName, cfg.MongoEnvDBName)
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("서버정보 서비스 생성 오류")
	}

	// 인증 서비스
	auth_service, err := auth.NewAuthService(server)
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("인증 서비스 생성 오류")
	}

	// 로그인 서비스
	login_service, err := login.NewLoginService(server, server.GetRouter())
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("로그인 서비스 생성 오류")
	}

	// 유저 서비스
	user_service, err := user.NewUserService(server, server.GetRouter())
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("유저 서비스 생성 오류")
	}

	if errs != nil {
		return errs
	}

	gameDBName, err := serverInfo_service.GetGameDBName()
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("게임 DB 이름 조회 오류")
	}

	sessionRepo, err := session.NewSessionRepository(server.GetContext(), server.GetMongoClient(), gameDBName)
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("세션 레포지토리 생성 오류")
	}

	userRepo, err := user.NewUserMongoRepository(server.GetMongoClient(), gameDBName)
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("유저 레포지토리 생성 오류")
	}

	// Register Handlers
	errs = nil

	if err := auth_service.SetRepositories(sessionRepo); err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("인증 서비스 레포지토리 설정 오류")
	}

	if err := login_service.SetHandlers(auth_service); err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("로그인 서비스 핸들러 설정 오류")
	}

	if err := login_service.SetRepositories(userRepo); err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("로그인 서비스 레포지토리 설정 오류")
	}

	if err := user_service.SetHandlers(auth_service); err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("유저 서비스 핸들러 설정 오류")
	}

	if err := user_service.SetRepositories(userRepo); err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("유저 서비스 레포지토리 설정 오류")
	}

	if errs != nil {
		return errs
	}

	// Add Services
	errs = nil
	for _, svc := range []service.Service{
		serverInfo_service,
		auth_service,
		user_service,
		login_service,
	} {
		if err := server.AddService(svc); err != nil {
			errs = errors.Join(errs, err)
			log.Error().Err(err).Msg("서비스 추가 오류")
		}
	}

	return errs
}

func run(ctx context.Context, cfg *config.WebServerConfig) error {
	opts := options.Client().ApplyURI(cfg.MongoUri)
	mongoClient, err := mongo.Connect(ctx, opts)
	log.Info().Msgf("MongoDB 연결을 시도 합니다. [uri:%v]", cfg.MongoUri)
	if err != nil {
		log.Err(err).Msgf("MongoDB 연결에 실패하였습니다. [uri:%v]", cfg.MongoUri)
		return err
	} else if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Err(err).Msgf("MongoDB 연결 실패하였습니다.. [uri:%v]", cfg.MongoUri)
		return err
	}

	log.Info().Msgf("MongoDB 연결 완료이 완료되었습니다. [uri:%v]", cfg.MongoUri)

	web_server, err := server.NewWebServer(
		ctx,
		cfg,
		mongoClient,
	)

	if err != nil {
		panic(err)
	}

	if err := setupServices(web_server, cfg); err != nil {
		panic(err)
	}

	if err := web_server.Init(); err != nil {
		panic(err)
	}

	serverErrCh := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("런타임 에러 발생: %v", r)
				serverErrCh <- err
			}
		}()

		if err := web_server.Start(); err != nil {
			serverErrCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		log.Info().Msg("서버가 강제종료 되었습니다.")

	case <-web_server.GetContext().Done():
		log.Info().Msg("서버가 종료되었습니다.")

	case err := <-serverErrCh:
		log.Err(err).Msg("치명적 서버 오류 발생가 발생하였습니다.")
		return err
	}

	return nil
}
