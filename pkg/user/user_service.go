package user

import (
	"MScannot206/pkg/user/mongo"
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

	s.userRepo, err = mongo.NewUserRepository(s.host.GetContext(), s.host.GetMongoClient(), dbName)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Start() error {
	if err := s.userRepo.Start(); err != nil {
		return err
	}

	return nil
}

func (s *UserService) Stop() error {
	if err := s.userRepo.Stop(); err != nil {
		return err
	}

	return nil
}
