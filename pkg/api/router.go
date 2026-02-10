package api

import (
	"MScannot206/pkg/api/batch"
	channel_api "MScannot206/pkg/api/channel"
	"MScannot206/pkg/api/login"
	"MScannot206/pkg/api/user"
	"MScannot206/shared/service"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type apiHandler interface {
	RegisterHandle(*http.ServeMux)
	Execute(ctx context.Context, api string, body json.RawMessage) (any, error)
	GetApiNames() []string
}

func SetupRoutes(host service.ServiceHost, r *http.ServeMux) error {
	if host == nil {
		return service.ErrServiceHostIsNil
	}

	if r == nil {
		return errors.New("router가 없습니다.")
	}

	apiManager := NewApiManager()

	var errs error

	// Batch
	batchHandler, err := batch.NewBatchHandler(host, apiManager)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	loginHandler, err := login.NewLoginHandler(host)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	userHandler, err := user.NewUserHandler(host)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	channelHandler, err := channel_api.NewChannelHandler(host)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	if errs != nil {
		return errs
	}

	// 배치 핸들러는 별도로 등록
	batchHandler.RegisterHandle(r)

	// bind
	for _, h := range []apiHandler{
		loginHandler,
		userHandler,
		channelHandler,
	} {
		// 핸들러 등록
		h.RegisterHandle(r)

		// API 호출기 등록
		for _, apiName := range h.GetApiNames() {
			apiManager.RegisterApiCaller(apiName, h)
		}
	}

	return errs
}
