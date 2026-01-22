package channel

import (
	channel_pkg "MScannot206/pkg/channel"
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

	channelService, err := service.GetService[*channel_pkg.ChannelService](host)
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

	channelService *channel_pkg.ChannelService
}

func (h *ChannelHandler) RegisterHandle(r *http.ServeMux) {
	r.HandleFunc("POST /api/v1/channel/create", h.HandleCreateChannel)
	r.HandleFunc("POST /api/v1/channel/renew", h.HandleRenewChannel)
	r.HandleFunc("GET /api/v1/channel/list", h.HandleListChannels)
}

func (h *ChannelHandler) HandleCreateChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req channel_pkg.AcquireChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	channel, err := h.channelService.Create(ctx, req.Id)
	if err != nil {
		log.Err(err).Msg("채널 생성 중 오류가 발생했습니다.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info().Str("channel_id", channel.Id).Msg("채널을 생성하였습니다.")

	channels, err := h.channelService.GetChannels(ctx)
	if err != nil {
		log.Err(err).Msg("채널들을 불러오는 중 오류가 발생했습니다.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var res channel_pkg.CreateChannelResponse
	res.Channels = channel_pkg.ToChannels(channels)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error().Err(err).Msg("채널 임대 응답 인코딩 중 오류가 발생했습니다.")
	}
}

func (h *ChannelHandler) HandleRenewChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req channel_pkg.RenewChannelRequest
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

	channels, err := h.channelService.GetChannels(ctx)
	if err != nil {
		log.Err(err).Msg("채널들을 불러오는 중 오류가 발생했습니다.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var res channel_pkg.RenewChannelResponse
	res.Channels = channel_pkg.ToChannels(channels)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error().Err(err).Msg("채널 갱신 응답 인코딩 중 오류가 발생했습니다.")
	}
}

func (h *ChannelHandler) HandleListChannels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	channels, err := h.channelService.GetChannels(ctx)
	if err != nil {
		log.Err(err).Msg("채널들을 불러오는 중 오류가 발생했습니다.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var res channel_pkg.ChannelListResponse
	res.Channels = channel_pkg.ToChannels(channels)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error().Err(err).Msg("채널 목록 응답 인코딩 중 오류가 발생했습니다.")
	}
}
