package login

// 로그인 요청 구조체
type LoginRequest struct {
	// 로그인 요청할 유저의 UID 목록
	Uids []string `json:"uids"`
}
