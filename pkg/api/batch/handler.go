package batch

import (
	"MScannot206/shared/service"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"
)

var ErrAPIManagerIsNil = errors.New("apiManager가 없습니다.")

type apiManager interface {
	ExecuteApi(ctx context.Context, wg *sync.WaitGroup, api string, body json.RawMessage) (<-chan *ApiResult, error)
}

func NewBatchHandler(
	host service.ServiceHost,
	am apiManager,
) (*BatchHandler, error) {
	if host == nil {
		return nil, service.ErrServiceHostIsNil
	}

	if am == nil {
		return nil, ErrAPIManagerIsNil
	}

	return &BatchHandler{
		host: host,
		am:   am,
	}, nil
}

type BatchHandler struct {
	host service.ServiceHost

	am apiManager
}

func (h *BatchHandler) RegisterHandle(r *http.ServeMux) {
	r.HandleFunc("POST /api/v1/batch", h.HandleBatch)
}

func (h *BatchHandler) HandleBatch(w http.ResponseWriter, r *http.Request) {
	var req []HttpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var dto []DataTransferObject
	for _, item := range req {
		dto = append(dto, item.Dto)
	}

	var res HttpResponse
	var wg sync.WaitGroup
	apiResults := make([]<-chan *ApiResult, 0, len(dto))
	resultCollector := make(chan *ApiResult, len(dto))

	for _, dto := range dto {
		ret, err := h.am.ExecuteApi(r.Context(), &wg, dto.Api, dto.Body)
		if err != nil {
			res.Dto = append(res.Dto, DataTransferObject{
				Api:       dto.Api,
				ErrorCode: BATCH_UNKNOWN_ERROR,
			})
			log.Err(err).Msg("API 호출기 실행 중 오류가 발생했습니다: " + dto.Api)
			continue
		}

		apiResults = append(apiResults, ret)
	}

	for _, ret := range apiResults {
		wg.Add(1)
		go func(resultChan <-chan *ApiResult) {
			defer wg.Done()
			resultCollector <- <-resultChan
		}(ret)
	}

	// 모든 작업이 완료될 때까지 대기
	go func() {
		wg.Wait()
		close(resultCollector)
	}()

	// 결과 수집
	for apiResult := range resultCollector {
		jsonBody, err := json.Marshal(apiResult.Body)
		if err != nil {
			res.Dto = append(res.Dto, DataTransferObject{
				Api:       apiResult.Api,
				ErrorCode: BATCH_UNKNOWN_ERROR,
			})
			log.Err(err).Msg("API 결과를 직렬화하는 중 오류가 발생했습니다.")
			continue
		}

		res.Dto = append(res.Dto, DataTransferObject{
			Api:       apiResult.Api,
			Body:      jsonBody,
			ErrorCode: apiResult.ErrorCode,
		})
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
