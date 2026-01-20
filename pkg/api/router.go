package api

import (
	"MScannot206/pkg/api/channel"
	"MScannot206/pkg/api/login"
	"MScannot206/pkg/api/user"
	"MScannot206/shared/service"
	"errors"
	"net/http"
)

type apiHandler interface {
	RegisterHandle(*http.ServeMux)
}

func SetupRoutes(host service.ServiceHost, r *http.ServeMux) error {
	if host == nil {
		return service.ErrServiceHostIsNil
	}

	if r == nil {
		return errors.New("router가 없습니다.")
	}

	var errs error

	loginHandler, err := login.NewLoginHandler(host)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	userHandler, err := user.NewUserHandler(host)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	channelHandler, err := channel.NewChannelHandler(host)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	// bind
	for _, h := range []apiHandler{
		loginHandler,
		userHandler,
		channelHandler,
	} {
		h.RegisterHandle(r)
	}

	return errs
}
