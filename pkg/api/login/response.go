package login

import "MScannot206/shared/entity"

type LoginSuccess struct {
	UserEntity *entity.User `json:"user_entity"`
	Token      string       `json:"token"`
}

type LoginFailure struct {
	Uid       string `json:"uid"`
	ErrorCode string `json:"error_code,omitempty"`
}

// 로그인 응답 구조체
type LoginResponse struct {
	Successes []*LoginSuccess `json:"successes"`
	Failures  []*LoginFailure `json:"failures"`
}
