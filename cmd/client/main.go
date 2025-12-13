package main

import (
	"MScannot206/pkg/manager"
	"MScannot206/pkg/testclient"
	"MScannot206/shared/client"
	"MScannot206/shared/config"
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
		AppName:   "TestClient",
		DebugMode: true,
	}

	var clientCfgPath string
	clientCfg := &config.WebClientConfig{
		Url:  "http://localhost",
		Port: 8080,
	}

	flag.StringVar(&logCfgPath, "logconfig", "", "로그 설정 파일 경로 지정")
	flag.StringVar(&clientCfgPath, "config", "", "클라이언트 설정 파일 경로 지정")
	flag.Parse()

	if logCfgPath != "" {
		if err := config.LoadYamlConfig(logCfgPath, logCfg); err != nil {
			defaultPath := filepath.Join(filepath.Dir(ex), "testclient_log_config.yaml")
			if _, err := os.Stat(defaultPath); err == nil {
				if err := config.LoadYamlConfig(defaultPath, logCfg); err != nil {
					log.Warn().Msg(err.Error())
				}
			}
		}
	} else {
		defaultPath := filepath.Join(filepath.Dir(ex), "testclient_log_config.yaml")
		if _, err := os.Stat(defaultPath); err == nil {
			if err := config.LoadYamlConfig(defaultPath, logCfg); err != nil {
				log.Warn().Msg(err.Error())
			}
		}
	}

	if err := manager.GetLogManager().Init(*logCfg); err != nil {
		println("로그 매니저 초기화 실패:", err)
	}
	defer manager.GetLogManager().Close()

	if clientCfgPath != "" {
		if err := config.LoadYamlConfig(clientCfgPath, clientCfg); err != nil {
			defaultPath := filepath.Join(filepath.Dir(ex), "testclient_config.yaml")
			if _, err := os.Stat(defaultPath); err == nil {
				if err := config.LoadYamlConfig(defaultPath, clientCfg); err != nil {
					log.Warn().Msg(err.Error())
				}
			}
		}
	} else {
		defaultPath := filepath.Join(filepath.Dir(ex), "testclient_config.yaml")
		if _, err := os.Stat(defaultPath); err == nil {
			if err := config.LoadYamlConfig(defaultPath, clientCfg); err != nil {
				log.Warn().Msg(err.Error())
			}
		}
	}

	client, err := client.NewWebClient(
		context.Background(),
		clientCfg,
	)

	if err != nil {
		panic(err)
	}

	// 코어 서비스
	core_service, err := testclient.NewCoreService(client)
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("코어 서비스 생성 오류")
	}

	// 로그인 서비스
	login_service, err := testclient.NewLoginService(client)
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("로그인 서비스 생성 오류")
	}

	if errs != nil {
		panic(errs)
	}

	errs = nil
	for _, svc := range []service.GenericService{
		core_service,
		login_service,
	} {
		if err := client.AddService(svc); err != nil {
			errs = errors.Join(errs, err)
			log.Error().Err(err).Msg("서비스 추가 오류")
		}
	}

	if errs != nil {
		panic(errs)
	}

	if err := client.Init(); err != nil {
		panic(err)
	}

	go func() {
		if err := client.Start(); err != nil {
			panic(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		log.Printf("클라이언트 강제 종료")

	case <-client.GetContext().Done():
		log.Printf("클라이언트 종료")
	}
}
