package app

import (
	"MScannot206/pkg/testclient/client"
	"MScannot206/pkg/testclient/config"
	"MScannot206/pkg/testclient/framework"
	"MScannot206/pkg/testclient/login"
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func CreateTestClient(ctx context.Context, cfg *config.ClientConfig) (framework.Client, error) {
	var errs error

	client, err := client.NewClient(
		ctx,
		cfg,
	)

	if err != nil {
		return nil, err
	}

	// 로그인 로직
	login_logic, err := login.NewLoginLogic(client)
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("로그인 서비스 생성 오류")
	}

	if errs != nil {
		return nil, errs
	}

	errs = nil
	for _, l := range []framework.Logic{
		login_logic,
	} {
		if err := client.AddLogic(l); err != nil {
			errs = errors.Join(errs, err)
			log.Error().Err(err).Msg("서비스 추가 오류")
		}
	}

	return client, nil
}

func Run(client framework.Client) error {
	if client == nil {
		return errors.New("client is null")
	}

	if err := client.Init(); err != nil {
		panic(err)
	}

	if err := RegisterCommands(client); err != nil {
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
		log.Info().Msg("클라이언트를 강제 종료합니다.")

	case <-client.GetContext().Done():
		log.Info().Msg("클라이언트 종료합니다.")
	}

	return nil
}

func RegisterCommands(client framework.Client) error {
	var errs error

	loginCommand, err := login.NewLoginCommand(client)
	if err != nil {
		errs = errors.Join(errs, err)
		log.Err(err)
	}

	if errs != nil {
		return errs
	}

	errs = nil
	for _, cmd := range []framework.ClientCommand{
		loginCommand,
	} {
		if err := client.AddCommand(cmd); err != nil {
			errs = errors.Join(errs, err)
			log.Err(err)
		}
	}

	return errs
}
