package api

import (
	"context"
	"sync"
)

func NewApiManager() *ApiManager {
	return &ApiManager{}
}

type ApiManager struct {
}

func (m *ApiManager) ExecuteApi(ctx context.Context, wg *sync.WaitGroup, api string, body string) (string, error) {
	// Implement the method logic here
	return "", nil
}
