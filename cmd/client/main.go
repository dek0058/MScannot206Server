package main

import (
	"MScannot206/pkg/logger"
	"MScannot206/pkg/testclient"
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

	if err := logger.GetLogManager().Init(*logCfg); err != nil {
		println("로그 매니저 초기화 실패:", err)
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

	client, err := testclient.CreateTestClient(context.Background(), clientCfg)

	if err != nil {
		panic(err)
	}

	if err := testclient.Run(client); err != nil {
		panic(err)
	}
}
