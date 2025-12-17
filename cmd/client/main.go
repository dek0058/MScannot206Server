package main

import (
	"MScannot206/pkg/logger"
	"MScannot206/pkg/testclient/app"
	testclient_config "MScannot206/pkg/testclient/config"
	"MScannot206/shared/config"
	"context"
	"flag"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func main() {
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
	clientCfg := &testclient_config.ClientConfig{
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

	if err := logger.GetLogManager().Init(*logCfg); err != nil {
		log.Err(err).Msg("로그 매니저 초기화에 실패하였습니다.")
	}
	defer logger.GetLogManager().Close()

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

	client, err := app.CreateTestClient(context.Background(), clientCfg)

	if err != nil {
		log.Err(err).Msg("테스트 클라이언트 생성 중 에러가 발생하였습니다.")
		panic(err)
	}

	if err := app.Run(client); err != nil {
		log.Err(err).Msg("테스트 클라이언트 실행 중 에러가 발생하였습니다.")
		panic(err)
	}
}
