package login

import (
	"MScannot206/pkg/login/mongodb"
	"MScannot206/shared/repository"
	"MScannot206/shared/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"
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
}

func (s *LoginService) Init(host service.ServiceHost) error {
	var err error

	// TODO: 외부에서 가져올 수 있도록 수정 필요
	dbName := "MStest"

	s.userRepo, err = mongodb.NewUserRepository(s.host.GetContext(), s.host.GetMongoClient(), dbName)
	if err != nil {
		return err
	}

	s.router.HandleFunc("/login", s.onLogin)

	return nil
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
				reason := LoginFailUid{
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

	// 로그인 프로세스
	for _, u := range users {
		response.SuccessUids = append(response.SuccessUids, u.Uid)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("onLogin JSON Encode Error: %v", err)
	}
}
