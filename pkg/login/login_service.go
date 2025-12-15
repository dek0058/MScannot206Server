package login

import (
	"MScannot206/pkg/auth"
	"MScannot206/pkg/user"
	"MScannot206/pkg/user/mongo"
	"MScannot206/shared/repository"
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
		return nil, errors.New("host is null")
	}

	if router == nil {
		return nil, errors.New("router is null")
	}

	return &LoginService{
		host:   host,
		router: router,
	}, nil
}

type LoginService struct {
	host   service.ServiceHost
	router *http.ServeMux

	userRepo repository.UserRepository

	authService *auth.AuthService
	userService *user.UserService
}

func (s *LoginService) Init() error {
	var errs error
	var err error

	s.router.HandleFunc("/login", s.onLogin)

	// TODO: 외부에서 가져올 수 있도록 수정 필요
	dbName := "MStest"

	s.userRepo, err = mongo.NewUserRepository(s.host.GetContext(), s.host.GetMongoClient(), dbName)
	if err != nil {
		log.Err(err)
		errs = errors.Join(errs, err)
	}

	s.authService, err = service.GetService[*auth.AuthService](s.host)
	if err != nil {
		log.Err(err)
		errs = errors.Join(errs, err)
	}

	s.userService, err = service.GetService[*user.UserService](s.host)
	if err != nil {
		log.Err(err)
		errs = errors.Join(errs, err)
	}

	return errs
}

func (s *LoginService) Start() error {
	var errs error

	if err := s.userRepo.Start(); err != nil {
		errs = errors.Join(errs, err)
	}

	if errs != nil {
		return errs
	}

	return nil
}

func (s *LoginService) Stop() error {
	var errs error

	if err := s.userRepo.Stop(); err != nil {
		errs = errors.Join(errs, err)
	}

	if errs != nil {
		return errs
	}

	return nil
}

func (s *LoginService) onLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var response LoginResponse

	users, newUids, err := s.userRepo.FindUserByUids(req.Uids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 신규 유저 생성
	if len(newUids) > 0 {
		newUsers, err := s.userRepo.InsertUserByUids(newUids)
		if err != nil {
			// 신규 유저는 로그인 불가
			log.Printf("신규 유저 생성 불가: %v", err)

			for _, uid := range newUids {
				reason := LoginFailure{
					Uid:       uid,
					ErrorCode: LOGIN_DB_WRITE_ERROR,
				}
				response.FailUids = append(response.FailUids, reason)
			}

			// TODO: 생성 불가능한 유저 uid 로그 추가
		} else {
			users = append(users, newUsers...)
		}
	}

	sessions, failureUsers, err := s.authService.CreateUserSessions(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, session := range sessions {
		success := LoginSuccess{
			Uid:   session.Uid,
			Token: session.Token,
		}
		response.SuccessUids = append(response.SuccessUids, success)
	}

	for _, u := range failureUsers {
		reason := LoginFailure{
			Uid:       u.Uid,
			ErrorCode: LOGIN_SESSION_CREATE_ERROR,
		}
		response.FailUids = append(response.FailUids, reason)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Err(err)
	}
}
