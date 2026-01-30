package api

import (
	"MScannot206/pkg/api/batch"
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

type ApiCaller interface {
	Execute(ctx context.Context, api string, body string) (any, error)
}

func NewApiManager() *ApiManager {
	return &ApiManager{
		callers: make(map[string]ApiCaller, 100),
	}
}

type ApiManager struct {
	callers map[string]ApiCaller
}

func (m *ApiManager) RegisterApiCaller(api string, caller ApiCaller) {
	m.callers[strings.ToLower(api)] = caller
}

func (m *ApiManager) ExecuteApi(ctx context.Context, wg *sync.WaitGroup, api string, body string) (<-chan *batch.ApiResult, error) {
	if api == "" {
		return nil, errors.New("빈 API 호출은 허용되지 않습니다")
	}

	apiName := api
	split := strings.Split(api, "?")
	if len(split) > 0 {
		apiName = split[0]
	}

	apiName = strings.ToLower(apiName)
	if m.callers[apiName] == nil {
		return nil, errors.New("등록되지 않은 API 호출기입니다: " + api)
	}

	caller := m.callers[apiName]
	resultChan := make(chan *batch.ApiResult, 1)

	wg.Go(func() {
		defer close(resultChan)

		if api == "" {
			resultChan <- &batch.ApiResult{
				Api:       api,
				ErrorCode: batch.BATCH_UNKNOWN_ERROR,
			}
			return
		}

		ret, err := caller.Execute(ctx, apiName, body)
		if err != nil {
			resultChan <- &batch.ApiResult{
				Api:       api,
				ErrorCode: batch.BATCH_UNKNOWN_ERROR,
			}
			log.Err(err).Str("api", api).Msg("API 실행 중 오류가 발생했습니다.")
			return
		}
		resultChan <- &batch.ApiResult{
			Api:  api,
			Body: ret,
		}
	})

	return resultChan, nil
}
