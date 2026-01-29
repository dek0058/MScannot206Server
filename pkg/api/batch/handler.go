package batch

import (
	"MScannot206/shared/service"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
)

var ErrAPIManagerIsNil = errors.New("apiManager가 없습니다.")

type apiManager interface {
	ExecuteApi(ctx context.Context, wg *sync.WaitGroup, api string, body string) (string, error)
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
	ctx := r.Context()

	var req HttpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var res HttpResponse

	var wg sync.WaitGroup
	for _, dto := range req.Dto {
		splits := strings.Split(dto.Api, "?")
		if len(splits) < 1 {
			continue
		}
		apiName := splits[0]
		body, err := h.am.ExecuteApi(ctx, &wg, apiName, dto.Body)
		if err != nil {
			res.Dto = append(res.Dto, DataTransferObject{
				Api:       dto.Api,
				ErrorCode: BATCH_UNKNOWN_ERROR,
				Body:      err.Error(),
			})
			continue
		}

		res.Dto = append(res.Dto, DataTransferObject{
			Api:       dto.Api,
			ErrorCode: "",
			Body:      body,
		})
	}

	// 모든 작업이 완료될 때까지 대기
	wg.Wait()

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
