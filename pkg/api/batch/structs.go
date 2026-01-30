package batch

type DataTransferObject struct {
	Api       string `json:"api"`
	Body      string `json:"body"`
	ErrorCode string `json:"error_code,omitempty"`
}

type ApiResult struct {
	Api       string
	Body      any
	ErrorCode string
}
