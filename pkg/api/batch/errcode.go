package batch

import "MScannot206/shared"

const BATCH_UNKNOWN_ERROR = "BATCH_UNKNOWN_ERROR"

func init() {
	shared.RegisterError(BATCH_UNKNOWN_ERROR, "알 수 없는 오류가 발생했습니다")
}
