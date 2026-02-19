package user

import (
	"MScannot206/shared/entity"
	"MScannot206/shared/table"
	"MScannot206/shared/table/view"
	"MScannot206/shared/types"
	"context"
	"errors"
)

func NewUserService(tableRepo *table.Repository) (*UserService, error) {
	return &UserService{}, nil
}

// 유저 서비스는 유저 관리 및 유저에 종속된 데이터를 관리하는 서비스입니다
type UserService struct {
	// 캐릭터 생성 테이블 뷰
	createCharacterView view.CreateCharacterView

	// 랜덤 서비스 핸들러
	randomServiceHandler RandomServiceHandler

	// 테이블 레포지토리
	tableRepo *table.Repository

	// 유저 DB 레포지토리
	userRepo *UserMongoRepository
}

func (s *UserService) Start(ctx context.Context) error {
	return nil
}

func (s *UserService) Stop(ctx context.Context) error {
	return nil
}

func (s *UserService) SetRepositories(
	tableRepo *table.Repository,
	userRepo *UserMongoRepository,
) error {
	var errs error

	s.userRepo = userRepo
	if userRepo == nil {
		errs = errors.Join(errs, ErrUserMongoRepositoryIsNil)
	}

	s.tableRepo = tableRepo
	if tableRepo == nil {
		errs = errors.Join(errs, table.ErrTableRepositoryIsNil)
	} else { // 테이블 뷰 생성

		// 캐릭터 생성 테이블 뷰 생성
		s.createCharacterView = view.NewCreateCharacterView(
			tableRepo.CreateCharacter,
			tableRepo.CreateCharacterHair,
			tableRepo.CreateCharacterFace,
			tableRepo.CreateCharacterCap,
			tableRepo.CreateCharacterCape,
			tableRepo.CreateCharacterCoat,
			tableRepo.CreateCharacterGlove,
			tableRepo.CreateCharacterLongCoat,
			tableRepo.CreateCharacterPants,
			tableRepo.CreateCharacterShoes,
			tableRepo.CreateCharacterFaceAcc,
			tableRepo.CreateCharacterEysAcc,
			tableRepo.CreateCharacterEarAcc,
			tableRepo.CreateCharacter1HWeapon,
			tableRepo.CreateCharacter2HWeapon,
			tableRepo.CreateCharacterSubWeapon,
			tableRepo.CreateCharacterEar,
			tableRepo.CreateCharacterSkin,
		)
	}

	return errs
}

func (s *UserService) SetHandlers(
	randomServiceHandler RandomServiceHandler,
) error {
	var errs error

	s.randomServiceHandler = randomServiceHandler
	if randomServiceHandler == nil {
		errs = errors.Join(errs, ErrRandomServiceHandlerIsNil)
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

func (s *UserService) CreateCharacterByUsers(ctx context.Context, createInfos []*UserCreateCharacter) (map[string]UserCreateCharacterResult, error) {
	if len(createInfos) == 0 {
		return map[string]UserCreateCharacterResult{}, nil
	}

	if s.randomServiceHandler == nil {
		return map[string]UserCreateCharacterResult{}, ErrRandomServiceHandlerIsNil
	}

	ret := make(map[string]UserCreateCharacterResult, len(createInfos))
	params := make([]*UserCreateCharacter, 0, len(createInfos))
	for _, info := range createInfos {
		result := UserCreateCharacterResult{}
		switch info.Gender {
		case types.GenderType_Male:
			result.Equips = s.createCharacterView.GetMale(s.randomServiceHandler.GetCharacterCreateSeed())
		case types.GenderType_Female:
			result.Equips = s.createCharacterView.GetFemale(s.randomServiceHandler.GetCharacterCreateSeed())
		default:
			result.ErrorCode = USER_CREATE_CHARACTER_GENDER_INVALID_ERROR
		}
		ret[info.Uid] = result
		if result.ErrorCode == "" {
			params = append(params, info)
		}
	}

	createdCharacters, failureUids, err := s.userRepo.CreateCharacters(ctx, params)
	if err != nil {
		return map[string]UserCreateCharacterResult{}, err
	}

	for uid, character := range createdCharacters {
		if result, ok := ret[uid]; ok {
			result.Character = character
			ret[uid] = result
		}
	}

	for uid, failureCode := range failureUids {
		if result, ok := ret[uid]; ok {
			result.ErrorCode = failureCode
			ret[uid] = result
		}
	}

	return ret, nil
}

func (s *UserService) DeleteCharactersByUsers(ctx context.Context, deleteInfos []*UserDeleteCharacter) ([]string, error) {
	if len(deleteInfos) == 0 {
		return []string{}, nil
	}
	return s.userRepo.DeleteCharacters(ctx, deleteInfos)
}
