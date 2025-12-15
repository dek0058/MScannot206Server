package login

import "MScannot206/pkg/manager"

const LOGIN_DB_WRITE_ERROR = "LOGIN_DB_WRITE_ERROR"
const LOGIN_LOGIN_UNABLE = "LOGIN_LOGIN_UNABLE"
const LOGIN_SESSION_CREATE_ERROR = "LOGIN_SESSION_CREATE_ERROR"

func init() {
	manager.RegisterError(LOGIN_DB_WRITE_ERROR, "데이터베이스에 쓰기 실패")
	manager.RegisterError(LOGIN_LOGIN_UNABLE, "로그인 불가 상태")
	manager.RegisterError(LOGIN_SESSION_CREATE_ERROR, "세션 생성 실패")
}
