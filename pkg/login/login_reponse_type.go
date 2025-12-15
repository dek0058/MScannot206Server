package login

type LoginSuccess struct {
	Uid   string `json:"uid"`
	Token string `json:"token"`
}

type LoginFailure struct {
	Uid       string `json:"uid"`
	ErrorCode string `json:"error_code"`
}

type LoginResponse struct {
	SuccessUids []LoginSuccess `json:"success_uids"`
	FailUids    []LoginFailure `json:"fail_uids"`
}
