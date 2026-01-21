package channel

import (
	"MScannot206/pkg/channel"
	"MScannot206/shared/service"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func NewChannelHandler(
	host service.ServiceHost,
) (*ChannelHandler, error) {
	if host == nil {
		return nil, service.ErrServiceHostIsNil
	}

	channelService, err := service.GetService[*channel.ChannelService](host)
	if err != nil {
		return nil, err
	}

	return &ChannelHandler{
		host: host,

		channelService: channelService,
	}, nil
}

type ChannelHandler struct {
	host service.ServiceHost

	channelService *channel.ChannelService
}

func (h *ChannelHandler) RegisterHandle(r *http.ServeMux) {
	r.HandleFunc("POST /api/v1/channel/acquire", h.HandleAcquireChannel)
	r.HandleFunc("POST /api/v1/channel/renew", h.HandleRenewChannel)
}

func (h *ChannelHandler) HandleAcquireChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req channel.AcquireChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	channel, err := h.channelService.Acquire(ctx, req.Id)
	if err != nil {
		log.Error().Err(err).Msg("채널 임대 중 오류가 발생했습니다.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info().Str("channel_id", channel.ID).Msg("채널을 임대하였습니다.")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(channel); err != nil {
		log.Error().Err(err).Msg("채널 임대 응답 인코딩 중 오류가 발생했습니다.")
	}
}

func (h *ChannelHandler) HandleRenewChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req channel.RenewChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	channel, err := h.channelService.Renew(ctx, req.Id)
	if err != nil {
		log.Error().Err(err).Msg("채널 갱신 중 오류가 발생했습니다.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if channel == nil {
		log.Warn().Str("channel_id", req.Id).Msg("채널을 찾을 수 없습니다.")
		http.Error(w, "채널을 찾을 수 없습니다.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(channel); err != nil {
		log.Error().Err(err).Msg("채널 갱신 응답 인코딩 중 오류가 발생했습니다.")
	}
}
