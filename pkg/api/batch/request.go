package batch

// HttpRequest에 대한 요청 구조체
type HttpRequest struct {
	// Dto는 실제 요청 데이터를 담고 있는 구조체
	Dto DataTransferObject `json:"dto"`
}
