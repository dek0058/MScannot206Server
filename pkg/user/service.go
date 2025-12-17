package user

import (
	"MScannot206/pkg/serverinfo"
	"MScannot206/shared/repository"
	"MScannot206/shared/service"
	"errors"

	"github.com/rs/zerolog/log"
)

func NewUserService(
	host service.ServiceHost,
) (*UserService, error) {
	if host == nil {
		return nil, errors.New("host is null")
	}
	return &UserService{
		host: host,
	}, nil
}

type UserService struct {
	host service.ServiceHost

	userRepo repository.UserRepository
}

func (s *UserService) GetPriority() int {
	return 0
}

func (s *UserService) Init() error {
	var errs error
	var err error
	var gameDBName string = ""

	serverInfoService, err := service.GetService[*serverinfo.ServerInfoService](s.host)
	if err != nil {
		log.Err(err)
		errs = errors.Join(errs, err)
	} else {
		srvInfo, err := serverInfoService.GetInfo()
		if err != nil {
			log.Err(err)
			errs = errors.Join(errs, err)
		} else {
			gameDBName = srvInfo.GameDBName
		}
	}

	s.userRepo, err = NewUserRepository(s.host.GetMongoClient(), gameDBName)
	if err != nil {
		return err
	}

	return errs
}

func (s *UserService) Start() error {
	return nil
}

func (s *UserService) Stop() error {
	return nil
}

func (s *UserService) ConnectUser(uid string, token string) error {

	return nil
}

func (s *UserService) DisconnectUser(uid string, token string) error {
	return nil
}
