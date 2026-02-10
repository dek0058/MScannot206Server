package batch

import "encoding/json"

// 배치 요청에서 각 API 호출에 대한 정보를 담고 있는 구조체입니다
type DataTransferObject struct {
	// 호출할 API의 경로를 나타냅니다
	Api string `json:"api"`

	// API 호출에 필요한 요청 본문을 JSON 형식으로 담고 있습니다
	Body json.RawMessage `json:"body"`

	// API 호출 중 발생한 오류 코드를 나타냅니다
	ErrorCode string `json:"error_code,omitempty"`
}

// 각 API 호출의 결과를 담고 있는 구조체입니다
type ApiResult struct {
	// 호출된 API의 경로를 나타냅니다
	Api string

	// API 호출의 응답 본문을 담고 있습니다
	Body any

	// API 호출 중 발생한 오류 코드를 나타냅니다
	ErrorCode string
}
