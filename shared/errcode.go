package shared

import "errors"

var errMaps map[string]string

func getErrorMaps() map[string]string {
	if errMaps == nil {
		errMaps = make(map[string]string)
	}
	return errMaps
}

func RegisterError(error_code string, message string) {
	em := getErrorMaps()
	em[error_code] = message
}

func ToError(error_code string) error {
	em := getErrorMaps()
	if msg, ok := em[error_code]; ok {
		return errors.New(msg)
	}
	return errors.New("unknown error: " + error_code)
}
