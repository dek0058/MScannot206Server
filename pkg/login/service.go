package login

import (
	"MScannot206/shared/entity"
	"context"
	"errors"
)

func NewLoginService() (*LoginService, error) {
	return &LoginService{}, nil
}

type LoginService struct {
	userRepoHandler UserRepositoryHandler
}

func (s *LoginService) Start(ctx context.Context) error {
	return nil
}

func (s *LoginService) Stop(ctx context.Context) error {
	return nil
}

func (s *LoginService) SetRepositories(
	userRepo UserRepositoryHandler,
) error {
	var errs error

	s.userRepoHandler = userRepo
	if userRepo == nil {
		errs = errors.Join(errs, ErrUserRepositoryHandlerIsNil)
	}

	return errs
}

func (s *LoginService) LoginUsers(ctx context.Context, uids []string) ([]*entity.User, error) {
	users, newUids, err := s.userRepoHandler.FindUserByUids(ctx, uids)
	if err != nil {
		return nil, err
	}

	// 신규 유저 생성
	if len(newUids) > 0 {
		newUsers, _, err := s.userRepoHandler.InsertUserByUids(ctx, newUids)
		if err != nil {
			return nil, err
		}
		users = append(users, newUsers...)
	}

	return users, nil
}
