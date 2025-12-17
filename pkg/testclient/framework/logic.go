package framework

import (
	"fmt"
)

type Logic interface {
	Init() error
	Start() error
	Stop() error
}

func GetLogic[T Logic](client Client) (T, error) {
	var nt T

	if client == nil {
		return nt, nil
	}

	for _, logic := range client.GetLogics() {
		if t, ok := logic.(T); ok {
			return t, nil
		}
	}

	return nt, fmt.Errorf("logic not found: %T", nt)
}
