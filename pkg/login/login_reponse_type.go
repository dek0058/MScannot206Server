package login

type LoginFailUid struct {
	Uid       string `json:"uid"`
	ErrorCode string `json:"error_code"`
}

type LoginResponse struct {
	SuccessUids []string       `json:"success_uids"`
	FailUids    []LoginFailUid `json:"fail_uids"`
}
