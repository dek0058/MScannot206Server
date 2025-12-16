package main

import (
	"MScannot206/pkg/auth"
	"MScannot206/pkg/logger"
	"MScannot206/pkg/login"
	"MScannot206/pkg/serverinfo"
	"MScannot206/pkg/user"
	"MScannot206/shared/config"
	"MScannot206/shared/server"
	"MScannot206/shared/service"
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	var errs error

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	var logCfgPath string
	logCfg := &config.LogConfig{
		AppName:   "server",
		DebugMode: true,
	}

	var serverCfgPath string
	serverCfg := &config.WebServerConfig{
		Port: 8080,

		ServerName: "DevServer",

		MongoUri:       "mongodb://localhost:27017/",
		MongoEnvDBName: "MSenv",
	}

	flag.StringVar(&logCfgPath, "logconfig", "", "로그 설정 파일 경로 지정")
	flag.StringVar(&serverCfgPath, "config", "", "서버 설정 파일 경로 지정")
	flag.Parse()

	if logCfgPath != "" {
		if err := config.LoadYamlConfig(logCfgPath, logCfg); err != nil {
			defaultPath := filepath.Join(filepath.Dir(ex), "server_log_config.yaml")
			if _, err := os.Stat(defaultPath); err == nil {
				if err := config.LoadYamlConfig(defaultPath, logCfg); err != nil {
					log.Warn().Msg(err.Error())
				}
			}
		}
	} else {
		defaultPath := filepath.Join(filepath.Dir(ex), "server_log_config.yaml")
		if _, err := os.Stat(defaultPath); err == nil {
			if err := config.LoadYamlConfig(defaultPath, logCfg); err != nil {
				log.Warn().Msg(err.Error())
			}
		}
	}

	if err := logger.GetLogManager().Init(*logCfg); err != nil {
		println("로그 매니저 초기화 실패:", err)
	}
	defer logger.GetLogManager().Close()

	if serverCfgPath != "" {
		if err := config.LoadYamlConfig(serverCfgPath, serverCfg); err != nil {
			defaultPath := filepath.Join(filepath.Dir(ex), "server_config.yaml")
			if _, err := os.Stat(defaultPath); err == nil {
				if err := config.LoadYamlConfig(defaultPath, serverCfg); err != nil {
					log.Warn().Msg(err.Error())
				}
			}
		}
	} else {
		defaultPath := filepath.Join(filepath.Dir(ex), "server_config.yaml")
		if _, err := os.Stat(defaultPath); err == nil {
			if err := config.LoadYamlConfig(defaultPath, serverCfg); err != nil {
				log.Warn().Msg(err.Error())
			}
		}
	}

	web_server, err := server.NewWebServer(
		context.Background(),
		serverCfg,
	)

	if err != nil {
		panic(err)
	}

	// 서버정보 서비스
	serverInfo_service, err := serverinfo.NewServerInfoService(web_server, serverCfg.ServerName, serverCfg.MongoEnvDBName)
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("서버정보 서비스 생성 오류")
	}

	// 인증 서비스
	auth_service, err := auth.NewAuthService(web_server)
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("인증 서비스 생성 오류")
	}

	// 로그인 서비스
	login_service, err := login.NewLoginService(web_server, web_server.GetRouter())
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("로그인 서비스 생성 오류")
	}

	// 유저 서비스
	user_service, err := user.NewUserService(web_server)
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("유저 서비스 생성 오류")
	}

	if errs != nil {
		panic(errs)
	}

	errs = nil
	for _, svc := range []service.GenericService{
		serverInfo_service,
		auth_service,
		user_service,
		login_service,
	} {
		if err := web_server.AddService(svc); err != nil {
			errs = errors.Join(errs, err)
			log.Error().Err(err).Msg("서비스 추가 오류")
		}
	}

	if errs != nil {
		panic(errs)
	}

	if err := web_server.Init(); err != nil {
		panic(err)
	}

	go func() {
		if err := web_server.Start(); err != nil {
			panic(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		log.Printf("서버 강제 종료")

	case <-web_server.GetContext().Done():
		log.Printf("서버 종료")
	}
}
