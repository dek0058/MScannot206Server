package login

import (
	"MScannot206/shared/entity"
	"MScannot206/shared/server"
	"MScannot206/shared/service"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
)

func NewLoginService(
	host service.ServiceHost,
	router *http.ServeMux,
) (*LoginService, error) {
	if host == nil {
		return nil, service.ErrServiceHostIsNil
	}

	if router == nil {
		return nil, server.ErrServeMuxIsNil
	}

	return &LoginService{
		host:   host,
		router: router,
	}, nil
}

type LoginService struct {
	host   service.ServiceHost
	router *http.ServeMux

	userRepoHandler UserRepositoryHandler

	authServiceHandler AuthServiceHandler
}

func (s *LoginService) Init() error {

	s.router.HandleFunc("POST /api/v1/login", s.onLogin)

	return nil
}

func (s *LoginService) Start() error {
	return nil
}

func (s *LoginService) Stop() error {
	return nil
}

func (s *LoginService) SetHandlers(
	authService AuthServiceHandler,
) error {
	var errs error

	s.authServiceHandler = authService
	if authService == nil {
		errs = errors.Join(errs, ErrAuthServiceHandlerIsNil)
	}

	return errs
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

func (s *LoginService) onLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res LoginResponse

	users, newUids, err := s.userRepoHandler.FindUserByUids(ctx, req.Uids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 신규 유저 생성
	if len(newUids) > 0 {
		newUsers, failedUids, err := s.userRepoHandler.InsertUserByUids(ctx, newUids)
		if err != nil {
			// 신규 유저는 로그인 불가
			log.Printf("신규 유저 생성 불가: %v", err)

			for _, uid := range newUids {
				reason := &LoginFailure{
					Uid:       uid,
					ErrorCode: LOGIN_DB_WRITE_ERROR,
				}
				res.FailUids = append(res.FailUids, reason)
			}
		} else {
			users = append(users, newUsers...)

			for _, uid := range failedUids {
				reason := &LoginFailure{
					Uid:       uid,
					ErrorCode: LOGIN_DB_WRITE_ERROR,
				}
				res.FailUids = append(res.FailUids, reason)
			}
		}
	}

	sessions, failureUsers, err := s.authServiceHandler.CreateUserSessions(ctx, users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usersByUid := make(map[string]*entity.User)
	for _, u := range users {
		usersByUid[u.Uid] = u
	}

	for _, session := range sessions {
		if u, ok := usersByUid[session.Uid]; ok {
			success := &LoginSuccess{
				UserEntity: u,
				Token:      session.Token,
			}
			res.SuccessUids = append(res.SuccessUids, success)
		} else {
			log.Warn().Msgf("세션은 존재하나 유저가 없음: %s", session.Uid)
		}
	}

	for _, u := range failureUsers {
		reason := &LoginFailure{
			Uid:       u.Uid,
			ErrorCode: LOGIN_SESSION_CREATE_ERROR,
		}
		res.FailUids = append(res.FailUids, reason)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Err(err)
	}
}
