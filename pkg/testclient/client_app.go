package testclient

import (
	"MScannot206/pkg/testclient/framework"
	"MScannot206/pkg/testclient/login"
	"MScannot206/shared/client"
	"MScannot206/shared/config"
	"MScannot206/shared/service"
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func CreateTestClient(ctx context.Context, cfg *config.WebClientConfig) (*client.WebClient, error) {
	var errs error

	client, err := client.NewWebClient(
		ctx,
		cfg,
	)

	if err != nil {
		return nil, err
	}

	// 코어 서비스
	core_service, err := framework.NewCoreService(client)
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("코어 서비스 생성 오류")
	}

	// 로그인 서비스
	login_service, err := login.NewLoginService(client)
	if err != nil {
		errs = errors.Join(errs, err)
		log.Error().Err(err).Msg("로그인 서비스 생성 오류")
	}

	if errs != nil {
		return nil, errs
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

	return client, nil
}

func Run(client *client.WebClient) error {
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
		log.Printf("클라이언트 강제 종료")

	case <-client.GetContext().Done():
		log.Printf("클라이언트 종료")
	}

	return nil
}

func RegisterCommands(client *client.WebClient) error {
	var errs error

	coreService, err := service.GetService[*framework.CoreService](client)
	if err != nil {
		return err
	}

	if err := coreService.AddCommand(login.NewLoginCommand(client)); err != nil {
		errs = errors.Join(errs, err)
	}

	return errs
}
