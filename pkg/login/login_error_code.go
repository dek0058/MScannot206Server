package login

import "MScannot206/pkg/manager"

const LOGIN_DB_WRITE_ERROR = "LOGIN_DB_WRITE_ERROR"
const LOGIN_LOGIN_UNABLE = "LOGIN_LOGIN_UNABLE"

func init() {
	manager.RegisterError(LOGIN_DB_WRITE_ERROR, "데이터베이스에 쓰기 실패")
	manager.RegisterError(LOGIN_LOGIN_UNABLE, "로그인 불가 상태")
}
