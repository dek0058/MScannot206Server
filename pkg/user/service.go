package user

import (
	"MScannot206/pkg/auth/session"
	"MScannot206/shared/def"
	"MScannot206/shared/entity"
	"MScannot206/shared/server"
	"MScannot206/shared/service"
	"MScannot206/shared/util"
	"encoding/json"
	"errors"
	"net/http"
	"unicode/utf8"
)

func NewUserService(
	host service.ServiceHost,
	router *http.ServeMux,
) (*UserService, error) {
	if host == nil {
		return nil, service.ErrServiceHostIsNil
	}

	if router == nil {
		return nil, server.ErrServeMuxIsNil
	}

	return &UserService{
		host:   host,
		router: router,
	}, nil
}

type UserService struct {
	host   service.ServiceHost
	router *http.ServeMux

	userRepo *UserMongoRepository

	authServiceHandler AuthServiceHandler
}

func (s *UserService) Init() error {
	s.router.HandleFunc("POST /api/v1/user/character/create", s.onCreateCharacter)
	s.router.HandleFunc("POST /api/v1/user/character/create/check_name", s.onCheckCharacterName)
	s.router.HandleFunc("POST /api/v1/user/character/delete", s.onDeleteCharacter)

	return nil
}

func (s *UserService) Start() error {
	return nil
}

func (s *UserService) Stop() error {
	return nil
}

func (s *UserService) SetHandlers(
	authServiceHandler AuthServiceHandler,
) error {
	var errs error

	s.authServiceHandler = authServiceHandler
	if authServiceHandler == nil {
		errs = errors.Join(errs, ErrAuthServiceHandlerIsNil)
	}

	return errs
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

func (s *UserService) onCreateCharacter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req CreateCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res CreateCharacterResponse

	requestCount := len(req.Requests)
	sessions := make([]*entity.UserSession, 0, requestCount)
	createInfos := make(map[string]*UserCreateCharacterInfo, requestCount)
	for _, entry := range req.Requests {
		// 잘못된 슬롯 검출
		if entry.Slot < 1 || entry.Slot > def.MaxCharacterSlot {
			res.Responses = append(res.Responses, &UserCreateCharacterResult{
				Uid:       entry.Uid,
				Slot:      entry.Slot,
				ErrorCode: USER_CHARACTER_SLOT_INVALID_ERROR,
			})
			continue
		}

		// 잘못된 이름 길이 검출
		nameLen := utf8.RuneCountInString(entry.Name)
		if nameLen < def.MinCharacterNameLength {
			res.Responses = append(res.Responses, &UserCreateCharacterResult{
				Uid:       entry.Uid,
				Slot:      entry.Slot,
				ErrorCode: USER_CREATE_CHARACTER_NAME_MIN_LENGTH_ERROR,
			})
			continue
		} else if nameLen > def.MaxCharacterNameLength {
			res.Responses = append(res.Responses, &UserCreateCharacterResult{
				Uid:       entry.Uid,
				Slot:      entry.Slot,
				ErrorCode: USER_CREATE_CHARACTER_NAME_MAX_LENGTH_ERROR,
			})
			continue
		} else if util.HasSpecialChar(entry.Name, s.host.GetLocale()) {
			res.Responses = append(res.Responses, &UserCreateCharacterResult{
				Uid:       entry.Uid,
				Slot:      entry.Slot,
				ErrorCode: USER_CREATE_CHARACTER_NAME_SPECIAL_CHAR_ERROR,
			})
			continue
		}

		s := &entity.UserSession{
			Uid:   entry.Uid,
			Token: entry.Token,
		}
		sessions = append(sessions, s)
		createInfos[entry.Uid] = entry
	}

	_, invalidUids, err := s.authServiceHandler.ValidateUserSessions(ctx, sessions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, uid := range invalidUids {
		delete(createInfos, uid)
		res.Responses = append(res.Responses, &UserCreateCharacterResult{
			Uid:       uid,
			ErrorCode: session.SESSION_TOKEN_INVALID_ERROR,
		})
	}

	userCharacters, err := s.userRepo.FindCharacters(ctx, func() []string {
		uids := make([]string, 0, len(createInfos))
		for uid := range createInfos {
			uids = append(uids, uid)
		}
		return uids
	}())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for uid, chars := range userCharacters {
		for _, ch := range chars {
			if ch.Slot == createInfos[uid].Slot {
				// 이미 해당 슬롯에 캐릭터가 존재함
				res.Responses = append(res.Responses, &UserCreateCharacterResult{
					Uid:       uid,
					Slot:      createInfos[uid].Slot,
					ErrorCode: USER_CHARACTER_SLOT_ALREADY_EXISTS_ERROR,
				})
				delete(createInfos, uid)
				break
			}
		}
	}

	userCreateInfos := make([]*UserCreateCharacterInfo, 0, len(createInfos))
	for _, info := range createInfos {
		userCreateInfos = append(userCreateInfos, info)
	}

	successUids, failedUids, err := s.userRepo.CreateCharacters(ctx, userCreateInfos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, uid := range successUids {
		res.Responses = append(res.Responses, &UserCreateCharacterResult{
			Uid:       uid,
			Slot:      createInfos[uid].Slot,
			ErrorCode: "",
		})
	}

	for _, uid := range failedUids {
		// 통상적으로 여기까지 오면 캐릭터 이름 중복 오류일 것이라 가정
		res.Responses = append(res.Responses, &UserCreateCharacterResult{
			Uid:       uid,
			Slot:      createInfos[uid].Slot,
			ErrorCode: USER_CHARACTER_NAME_ALREADY_EXISTS_ERROR,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *UserService) onCheckCharacterName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req CheckCharacterNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res CheckCharacterNameResponse

	requestCount := len(req.Requests)
	sessions := make([]*entity.UserSession, 0, requestCount)
	nameCheckInfos := make(map[string]*UserNameCheckInfo, requestCount)
	for _, entry := range req.Requests {
		// 잘못된 이름 길이 검출
		nameLen := utf8.RuneCountInString(entry.Name)
		if nameLen < def.MinCharacterNameLength {
			res.Responses = append(res.Responses, &UserNameCheckResult{
				Uid:       entry.Uid,
				Available: false,
				ErrorCode: USER_CREATE_CHARACTER_NAME_MIN_LENGTH_ERROR,
			})
			continue
		} else if nameLen > def.MaxCharacterNameLength {
			res.Responses = append(res.Responses, &UserNameCheckResult{
				Uid:       entry.Uid,
				Available: false,
				ErrorCode: USER_CREATE_CHARACTER_NAME_MAX_LENGTH_ERROR,
			})
			continue
		} else if util.HasSpecialChar(entry.Name, s.host.GetLocale()) {
			res.Responses = append(res.Responses, &UserNameCheckResult{
				Uid:       entry.Uid,
				Available: false,
				ErrorCode: USER_CREATE_CHARACTER_NAME_SPECIAL_CHAR_ERROR,
			})
			continue
		}

		s := &entity.UserSession{
			Uid:   entry.Uid,
			Token: entry.Token,
		}
		sessions = append(sessions, s)
		nameCheckInfos[entry.Uid] = entry
	}

	_, invalidUids, err := s.authServiceHandler.ValidateUserSessions(ctx, sessions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, uid := range invalidUids {
		delete(nameCheckInfos, uid)
		res.Responses = append(res.Responses, &UserNameCheckResult{
			Uid:       uid,
			Available: false,
			ErrorCode: session.SESSION_TOKEN_INVALID_ERROR,
		})
	}

	chNames := make([]string, 0, len(nameCheckInfos))
	for _, info := range nameCheckInfos {
		chNames = append(chNames, info.Name)
	}

	existsMap, err := s.userRepo.ExistsCharacterNames(ctx, chNames)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, info := range nameCheckInfos {
		exists, ok := existsMap[info.Name]

		var ret *UserNameCheckResult
		if ok && exists {
			ret = &UserNameCheckResult{
				Uid:       info.Uid,
				Available: false,
				ErrorCode: USER_CHARACTER_NAME_ALREADY_EXISTS_ERROR,
			}
		} else {
			ret = &UserNameCheckResult{
				Uid:       info.Uid,
				Available: true,
				ErrorCode: "",
			}
		}

		res.Responses = append(res.Responses, ret)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *UserService) onDeleteCharacter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res DeleteCharacterResponse

	requestCount := len(req.Requests)
	sessions := make([]*entity.UserSession, 0, requestCount)
	deleteInfos := make(map[string]*UserDeleteCharacterInfo, requestCount)
	for _, entry := range req.Requests {
		// 잘못된 슬롯 검출
		if entry.Slot < 1 || entry.Slot > def.MaxCharacterSlot {
			res.Responses = append(res.Responses, &UserDeleteCharacterResult{
				Uid:       entry.Uid,
				ErrorCode: USER_DELETE_CHARACTER_SLOT_INVALID_ERROR,
			})
			continue
		}

		s := &entity.UserSession{
			Uid:   entry.Uid,
			Token: entry.Token,
		}
		sessions = append(sessions, s)
		deleteInfos[entry.Uid] = entry
	}

	_, invalidUids, err := s.authServiceHandler.ValidateUserSessions(ctx, sessions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, uid := range invalidUids {
		delete(deleteInfos, uid)
		res.Responses = append(res.Responses, &UserDeleteCharacterResult{
			Uid:       uid,
			ErrorCode: session.SESSION_TOKEN_INVALID_ERROR,
		})
	}

	userCharacters, err := s.userRepo.FindCharacters(ctx, func() []string {
		uids := make([]string, 0, len(deleteInfos))
		for uid := range deleteInfos {
			uids = append(uids, uid)
		}
		return uids
	}())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for uid, chars := range userCharacters {
		found := false
		for _, ch := range chars {
			if ch.Slot == deleteInfos[uid].Slot {
				deleteInfos[uid].Name = ch.Name
				found = true
				break
			}
		}
		if !found {
			delete(deleteInfos, uid)
			res.Responses = append(res.Responses, &UserDeleteCharacterResult{
				Uid:       uid,
				ErrorCode: USER_DELETE_CHARACTER_SLOT_NOT_FOUND_ERROR,
			})
		}
	}

	userDeleteInfos := make([]*UserDeleteCharacterInfo, 0, len(deleteInfos))
	for _, info := range deleteInfos {
		userDeleteInfos = append(userDeleteInfos, info)
	}
	successUids, failedUids, err := s.userRepo.DeleteCharacters(ctx, userDeleteInfos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, uid := range successUids {
		res.Responses = append(res.Responses, &UserDeleteCharacterResult{
			Uid:       uid,
			ErrorCode: "",
		})
	}

	for _, uid := range failedUids {
		res.Responses = append(res.Responses, &UserDeleteCharacterResult{
			Uid:       uid,
			ErrorCode: USER_DELETE_CHARACTER_DB_WRITE_ERROR,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
