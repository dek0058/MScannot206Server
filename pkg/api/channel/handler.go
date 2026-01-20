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
	r.HandleFunc("POST /api/v1/channel/start", h.HandleStartChannel)
	r.HandleFunc("POST /api/v1/channel/stop", h.HandleStopChannel)
	r.HandleFunc("GET /api/v1/channel/list", h.HandleListChannel)
}

func (h *ChannelHandler) HandleStartChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req channel.StartChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.channelService.AddChannel(ctx, req.Id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info().Str("channel_id", req.Id).Msg("채널이 시작되었습니다.")
	w.WriteHeader(http.StatusOK)
}

func (h *ChannelHandler) HandleStopChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req channel.StopChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.channelService.RemoveChannel(ctx, req.Id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info().Str("channel_id", req.Id).Msg("채널이 종료되었습니다.")
	w.WriteHeader(http.StatusOK)
}

func (h *ChannelHandler) HandleListChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	channels, err := h.channelService.GetChannels(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var res channel.FetchChannelsResponse
	res.Channels = channels

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
