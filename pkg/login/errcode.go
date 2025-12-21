package login

import "MScannot206/shared"

const LOGIN_UNKNOWN_ERROR = "LOGIN_UNKNOWN_ERROR"
const LOGIN_DB_WRITE_ERROR = "LOGIN_DB_WRITE_ERROR"
const LOGIN_UNABLE = "LOGIN_UNABLE"
const LOGIN_SESSION_CREATE_ERROR = "LOGIN_SESSION_CREATE_ERROR"

func init() {
	shared.RegisterError(LOGIN_UNKNOWN_ERROR, "알 수 없는 오류가 발생했습니다")
	shared.RegisterError(LOGIN_DB_WRITE_ERROR, "데이터베이스에 쓰기 오류가 발생했습니다")
	shared.RegisterError(LOGIN_UNABLE, "로그인이 불가능합니다")
	shared.RegisterError(LOGIN_SESSION_CREATE_ERROR, "세션 생성에 실패했습니다")
}
