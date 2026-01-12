package user

import (
	"MScannot206/shared/entity"
	"context"
	"errors"
)

func NewUserService() (*UserService, error) {
	return &UserService{}, nil
}

type UserService struct {
	userRepo *UserMongoRepository
}

func (s *UserService) Init() error {
	return nil
}

func (s *UserService) Start() error {
	return nil
}

func (s *UserService) Stop() error {
	return nil
}

func (s *UserService) SetRepositories(
	userRepo *UserMongoRepository,
) error {
	var errs error

	s.userRepo = userRepo
	if userRepo == nil {
		errs = errors.Join(errs, ErrUserMongoRepositoryIsNil)
	}

	return errs
}

func (s *UserService) FindCharactersByUids(ctx context.Context, uids []string) (map[string][]*entity.Character, error) {
	if len(uids) == 0 {
		return map[string][]*entity.Character{}, nil
	}
	return s.userRepo.FindCharacters(ctx, uids)
}

func (s *UserService) FindCharacterNames(ctx context.Context, names []string) (map[string]bool, error) {
	return s.userRepo.ExistsCharacterNames(ctx, names)
}

func (s *UserService) CreateCharacterByUsers(ctx context.Context, createInfos []*UserCreateCharacter) (map[string]*entity.Character, map[string]string, error) {
	if len(createInfos) == 0 {
		return map[string]*entity.Character{}, map[string]string{}, nil
	}
	return s.userRepo.CreateCharacters(ctx, createInfos)
}

func (s *UserService) DeleteCharactersByUsers(ctx context.Context, deleteInfos []*UserDeleteCharacter) ([]string, error) {
	if len(deleteInfos) == 0 {
		return []string{}, nil
	}
	return s.userRepo.DeleteCharacters(ctx, deleteInfos)
}
