package user

import (
	"MScannot206/shared/repository"
	"MScannot206/shared/service"
	"errors"
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

func (s *UserService) Init() error {
	var err error

	dbName := "MStest"

	s.userRepo, err = NewUserRepository(s.host.GetMongoClient(), dbName)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Start() error {
	return nil
}

func (s *UserService) Stop() error {
	return nil
}
