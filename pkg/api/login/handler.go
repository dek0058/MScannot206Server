package login

import (
	"MScannot206/pkg/auth"
	"MScannot206/pkg/login"
	"MScannot206/shared/entity"
	"MScannot206/shared/service"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func NewLoginHandler(
	host service.ServiceHost,
) (*LoginHandler, error) {
	if host == nil {
		return nil, service.ErrServiceHostIsNil
	}

	loginService, err := service.GetService[*login.LoginService](host)
	if err != nil {
		return nil, err
	}

	authService, err := service.GetService[*auth.AuthService](host)
	if err != nil {
		return nil, err
	}

	return &LoginHandler{
		host: host,

		loginService: loginService,
		authService:  authService,
	}, nil
}

type LoginHandler struct {
	host service.ServiceHost

	loginService *login.LoginService
	authService  *auth.AuthService
}

func (h *LoginHandler) RegisterHandle(r *http.ServeMux) {
	r.HandleFunc("POST /api/v1/login", h.HandleLogin)
}

func (h *LoginHandler) GetApiNames() []string {
	return []string{
		"login",
	}
}

func (h *LoginHandler) Execute(ctx context.Context, api string, body json.RawMessage) (any, error) {
	switch api {
	case "login":
		return h.login(ctx, body)

	default:
		return nil, errors.New("알 수 없는 API 호출입니다: " + api)
	}
}

func (h *LoginHandler) login(ctx context.Context, body json.RawMessage) (any, error) {
	var req LoginRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	failureUids := make(map[string]struct{})
	for _, uid := range req.Uids {
		failureUids[uid] = struct{}{}
	}

	var res LoginResponse

	// 로그인 처리
	users, err := h.loginService.LoginUsers(ctx, req.Uids)
	if err != nil {
		return nil, err
	}

	// 세션 생성
	sessions, failureUsers, err := h.authService.CreateUserSessions(ctx, users)
	if err != nil {
		return nil, err
	}

	// 로그인에 성공한 유저에 대한 처리 및 실패한 유저에 대한 처리
	loggedinUsers := make(map[string]*entity.User, len(users))
	for _, u := range users {
		loggedinUsers[u.Uid] = u
		delete(failureUids, u.Uid)
	}

	// 로그인에 실패한 유저의 경우 DB 쓰기 오류로 간주, 신규 유저는 새로 생성하기 때문
	for uid := range failureUids {
		reason := &LoginFailure{
			Uid:       uid,
			ErrorCode: login.LOGIN_DB_WRITE_ERROR,
		}
		res.Failures = append(res.Failures, reason)
	}

	// 세션 생성에 실패한 유저 추가
	for _, u := range failureUsers {
		delete(failureUids, u.Uid)

		// 이유 추가
		reason := &LoginFailure{
			Uid:       u.Uid,
			ErrorCode: login.LOGIN_SESSION_CREATE_ERROR,
		}
		res.Failures = append(res.Failures, reason)
	}

	// 최종적으로 성공한 유저 리스폰 정보 생성
	for _, s := range sessions {
		success := &LoginSuccess{
			UserEntity: loggedinUsers[s.Uid],
			Token:      s.Token,
		}
		res.Successes = append(res.Successes, success)
	}

	return &res, nil
}

func (h *LoginHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ret, err := h.login(ctx, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, ok := ret.(*LoginResponse)
	if !ok {
		http.Error(w, "응답 변환에 실패했습니다.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
