package user

import (
	"MScannot206/pkg/auth"
	"MScannot206/pkg/auth/session"
	"MScannot206/pkg/user"
	"MScannot206/shared/entity"
	"MScannot206/shared/service"
	"encoding/json"
	"net/http"
)

func NewUserHandler(
	host service.ServiceHost,
) (*UserHandler, error) {
	if host == nil {
		return nil, service.ErrServiceHostIsNil
	}

	authService, err := service.GetService[*auth.AuthService](host)
	if err != nil {
		return nil, err
	}

	userService, err := service.GetService[*user.UserService](host)
	if err != nil {
		return nil, err
	}

	return &UserHandler{
		host:        host,
		authService: authService,
		userService: userService,
	}, nil
}

type UserHandler struct {
	host service.ServiceHost

	authService *auth.AuthService
	userService *user.UserService
}

func (h *UserHandler) RegisterHandle(r *http.ServeMux) {
	r.HandleFunc("POST /api/v1/user/character/create", h.onCreateCharacter)
	r.HandleFunc("POST /api/v1/user/character/create/check_name", h.onCheckCharacterName)
	r.HandleFunc("POST /api/v1/user/character/delete", h.onDeleteCharacter)
}

// 캐릭터 생성 핸들러
func (h *UserHandler) onCreateCharacter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req user.CreateCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestCount := len(req.Requests)
	sessions := make([]*entity.UserSession, 0, requestCount)
	requests := make(map[string]*user.UserCreateCharacter, requestCount)
	var res user.CreateCharacterResponse

	for _, entry := range req.Requests {
		errCode := ""

		// 캐릭터 슬롯 유효성 검사
		if user.IsInvalidCharacterSlot(entry.Slot) {
			errCode = user.USER_CHARACTER_SLOT_INVALID_ERROR
		} else {
			// 캐릭터 이름 유효성 검사
			errCode = user.ValidateCharacterName(entry.Name, h.host.GetLocale())
		}

		// 오류가 있을 경우 다음 요청으로 넘어감
		if errCode != "" {
			res.Responses = append(res.Responses, &user.UserCreateCharacterResult{
				Uid:       entry.Uid,
				ErrorCode: errCode,
			})
			continue
		}

		sessions = append(sessions, &entity.UserSession{
			Uid:   entry.Uid,
			Token: entry.Token,
		})

		requests[entry.Uid] = &user.UserCreateCharacter{
			Uid:  entry.Uid,
			Slot: entry.Slot,
			Name: entry.Name,
		}
	}

	_, invalidUids, err := h.authService.ValidateUserSessions(ctx, sessions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, uid := range invalidUids {
		delete(requests, uid)
		res.Responses = append(res.Responses, &user.UserCreateCharacterResult{
			Uid:       uid,
			ErrorCode: session.SESSION_TOKEN_INVALID_ERROR,
		})
	}

	userCharacters, err := h.userService.FindCharactersByUids(ctx, func() []string {
		uids := make([]string, 0, len(requests))
		for uid := range requests {
			uids = append(uids, uid)
		}
		return uids
	}())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for uid, characters := range userCharacters {
		for _, ch := range characters {
			if ch.Slot != requests[uid].Slot {
				continue
			}

			res.Responses = append(res.Responses, &user.UserCreateCharacterResult{
				Uid:       uid,
				ErrorCode: user.USER_CHARACTER_SLOT_ALREADY_EXISTS_ERROR,
			})
			delete(requests, uid)
			break
		}
	}

	createdCharacters, failureUids, err := h.userService.CreateCharacterByUsers(ctx, func() []*user.UserCreateCharacter {
		createInfos := make([]*user.UserCreateCharacter, 0, len(requests))
		for _, info := range requests {
			createInfos = append(createInfos, info)
		}
		return createInfos
	}())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for uid := range requests {
		character, ok := createdCharacters[uid]
		// 생성된 캐릭터가 없을 경우
		if !ok {
			errorCode := user.USER_CREATE_CHARACTER_DB_WRITE_ERROR
			if _, ok := failureUids[uid]; ok {
				errorCode = failureUids[uid]
			}

			res.Responses = append(res.Responses, &user.UserCreateCharacterResult{
				Uid:       uid,
				ErrorCode: errorCode,
			})
			continue
		}

		res.Responses = append(res.Responses, &user.UserCreateCharacterResult{
			Uid:       uid,
			Character: character,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 캐릭터 이름 중복 확인 핸들러
func (h *UserHandler) onCheckCharacterName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req user.CheckCharacterNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestCount := len(req.Requests)
	sessions := make([]*entity.UserSession, 0, requestCount)
	requests := make(map[string]*user.UserCreateCharacter, requestCount)
	var res user.CheckCharacterNameResponse

	for _, entry := range req.Requests {
		// 캐릭터 이름 유효성 검사
		if errCode := user.ValidateCharacterName(entry.Name, h.host.GetLocale()); errCode != "" {
			res.Responses = append(res.Responses, &user.UserNameCheckResult{
				Uid:       entry.Uid,
				ErrorCode: errCode,
			})
			continue
		}

		sessions = append(sessions, &entity.UserSession{
			Uid:   entry.Uid,
			Token: entry.Token,
		})

		requests[entry.Uid] = &user.UserCreateCharacter{
			Uid:  entry.Uid,
			Name: entry.Name,
		}
	}

	_, invalidUids, err := h.authService.ValidateUserSessions(ctx, sessions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, uid := range invalidUids {
		delete(requests, uid)
		res.Responses = append(res.Responses, &user.UserNameCheckResult{
			Uid:       uid,
			ErrorCode: session.SESSION_TOKEN_INVALID_ERROR,
		})
	}

	existingNames, err := h.userService.FindCharacterNames(ctx, func() []string {
		names := make([]string, 0, len(requests))
		for _, info := range requests {
			names = append(names, info.Name)
		}
		return names
	}())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for uid, info := range requests {
		if exists, ok := existingNames[info.Name]; ok && exists {
			res.Responses = append(res.Responses, &user.UserNameCheckResult{
				Uid:       uid,
				ErrorCode: user.USER_CHARACTER_NAME_ALREADY_EXISTS_ERROR,
			})
		} else {
			res.Responses = append(res.Responses, &user.UserNameCheckResult{
				Uid: uid,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 캐릭터 삭제 핸들러
func (h *UserHandler) onDeleteCharacter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req user.DeleteCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestCount := len(req.Requests)
	sessions := make([]*entity.UserSession, 0, requestCount)
	requests := make(map[string]*user.UserDeleteCharacterInfo, requestCount)

	var res user.DeleteCharacterResponse

	for _, entry := range req.Requests {
		if user.IsInvalidCharacterSlot(entry.Slot) {
			res.Responses = append(res.Responses, &user.UserDeleteCharacterResult{
				Uid:       entry.Uid,
				ErrorCode: user.USER_CHARACTER_SLOT_INVALID_ERROR,
			})
			continue
		}

		sessions = append(sessions, &entity.UserSession{
			Uid:   entry.Uid,
			Token: entry.Token,
		})

		requests[entry.Uid] = &user.UserDeleteCharacterInfo{
			Uid:   entry.Uid,
			Token: entry.Token,
			Slot:  entry.Slot,
		}
	}

	_, invalidUids, err := h.authService.ValidateUserSessions(ctx, sessions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, uid := range invalidUids {
		delete(requests, uid)
		res.Responses = append(res.Responses, &user.UserDeleteCharacterResult{
			Uid:       uid,
			ErrorCode: session.SESSION_TOKEN_INVALID_ERROR,
		})
	}

	userCharacters, err := h.userService.FindCharactersByUids(ctx, func() []string {
		uids := make([]string, 0, len(requests))
		for uid := range requests {
			uids = append(uids, uid)
		}
		return uids
	}())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userDeleteCharacters := make([]*user.UserDeleteCharacter, 0, len(requests))
	for uid, characters := range userCharacters {
		var foundCharacter *entity.Character
		for _, ch := range characters {
			if ch.Slot != requests[uid].Slot {
				continue
			}
			foundCharacter = ch
			break
		}

		if foundCharacter == nil {
			res.Responses = append(res.Responses, &user.UserDeleteCharacterResult{
				Uid:       uid,
				ErrorCode: user.USER_DELETE_CHARACTER_SLOT_NOT_FOUND_ERROR,
			})
			delete(requests, uid)
		} else {
			// 삭제 목록에 해당 슬롯의 캐릭터 이름 주입
			userDeleteCharacters = append(userDeleteCharacters, &user.UserDeleteCharacter{
				Uid:  uid,
				Slot: foundCharacter.Slot,
				Name: foundCharacter.Name,
			})
		}
	}

	// 캐릭터 삭제 처리
	successUids, err := h.userService.DeleteCharactersByUsers(ctx, userDeleteCharacters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, uid := range successUids {
		if _, ok := requests[uid]; ok {
			res.Responses = append(res.Responses, &user.UserDeleteCharacterResult{
				Uid: uid,
			})
		} else {
			res.Responses = append(res.Responses, &user.UserDeleteCharacterResult{
				Uid:       uid,
				ErrorCode: user.USER_DELETE_CHARACTER_DB_WRITE_ERROR,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
