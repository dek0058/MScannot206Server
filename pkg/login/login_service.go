package login

import (
	"MScannot206/pkg/login/mongodb"
	"MScannot206/shared/repository"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewLoginService(
	ctx context.Context,
	router *http.ServeMux,
) (*LoginService, error) {
	if ctx == nil {
		return nil, errors.New("context is null")
	}

	if router == nil {
		return nil, errors.New("router is null")
	}

	return &LoginService{
		ctx:    ctx,
		router: router,
	}, nil
}

type LoginService struct {
	ctx    context.Context
	router *http.ServeMux
	client *mongo.Client

	userRepo repository.UserRepository
}

func (s *LoginService) Init() error {
	var err error
	s.userRepo, err = mongodb.NewUserRepository(s.client)
	if err != nil {
		return err
	}

	s.router.HandleFunc("/login", s.onLogin)

	return nil
}

func (s *LoginService) Start() error {
	return nil
}

func (s *LoginService) Stop() error {
	return nil
}

func (s *LoginService) onLogin(w http.ResponseWriter, r *http.Request) {

	// HTTP POST만 허용
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	uid := req.Uid

	user, err := s.userRepo.FindUserByUID(uid)
	if err != nil {
		http.Error(w, "User Not Found", http.StatusInternalServerError)
		return
	}

	_ = user

	// TODO: 접속중인 계정인지 확인
	println("User logged in:", uid)

}

/*

// 로그인 요청을 위한 구조체
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 로그인 응답을 위한 구조체
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"` // 로그인이 성공하면 토큰을 포함
}

func (s *LoginService) loginHandler(w http.ResponseWriter, r *http.Request) {
	println("Login Handler!")

	// 1. HTTP 메서드 확인 (POST만 허용)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. 요청 본문(Body) 디코딩
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// 3. 사용자 인증 로직 (⭐핵심⭐)
	// 이 부분에서 실제 데이터베이스를 조회하여 사용자 이름과 비밀번호를 검증해야 합니다.
	if req.Username == "testuser" && req.Password == "password123" {
		// 인증 성공
		response := LoginResponse{
			Success: true,
			Message: "Login successful",
			Token:   "dummy-jwt-token-12345", // 실제로는 JWT 등을 생성
		}

		// 4. 응답 전송
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		// 인증 실패
		response := LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized) // 401 Unathorized
		json.NewEncoder(w).Encode(response)
	}
}


*/
