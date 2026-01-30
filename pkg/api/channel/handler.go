package channel

import (
	channel_pkg "MScannot206/pkg/channel"
	"MScannot206/shared/service"
	"context"
	"encoding/json"
	"errors"
	"io"
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

func (h *ChannelHandler) GetApiNames() []string {
	return []string{
		"channel/create",
		"channel/renew",
		"channel/list",
	}
}

func (h *ChannelHandler) Execute(ctx context.Context, api string, body string) (any, error) {
	switch api {
	default:
		return nil, errors.New("알 수 없는 API 호출입니다: " + api)
	}
}

func (h *ChannelHandler) createChannel(ctx context.Context, body string) (any, error) {
	var req channel_pkg.AcquireChannelRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		log.Error().Err(err).Msg("채널 생성 요청 파싱 중 오류가 발생했습니다.")
		return nil, err
	}

	channel, err := h.channelService.Create(ctx, req.Id)
	if err != nil {
		log.Err(err).Msg("채널 생성 중 오류가 발생했습니다.")
		return nil, err
	}

	log.Info().Str("channel_id", channel.Id).Msg("채널을 생성하였습니다.")

	channels, err := h.channelService.GetChannels(ctx)
	if err != nil {
		log.Err(err).Msg("채널들을 불러오는 중 오류가 발생했습니다.")
		return nil, err
	}

	var res channel_pkg.CreateChannelResponse
	res.Channels = channel_pkg.ToChannels(channels)

	return res, nil
}

func (h *ChannelHandler) renewChannel(ctx context.Context, body string) (any, error) {
	var req channel_pkg.RenewChannelRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		log.Error().Err(err).Msg("채널 갱신 요청 파싱 중 오류가 발생했습니다.")
		return nil, err
	}

	channel, err := h.channelService.Renew(ctx, req.Id)
	if err != nil {
		log.Error().Err(err).Msg("채널 갱신 중 오류가 발생했습니다.")
		return nil, err
	}

	if channel == nil {
		return nil, errors.New("채널을 찾을 수 없습니다.")
	}

	channels, err := h.channelService.GetChannels(ctx)
	if err != nil {
		log.Err(err).Msg("채널들을 불러오는 중 오류가 발생했습니다.")
		return nil, err
	}

	var res channel_pkg.RenewChannelResponse
	res.Channels = channel_pkg.ToChannels(channels)

	return res, nil
}

func (h *ChannelHandler) listChannels(ctx context.Context, body string) (any, error) {
	channels, err := h.channelService.GetChannels(ctx)
	if err != nil {
		log.Err(err).Msg("채널들을 불러오는 중 오류가 발생했습니다.")
		return nil, err
	}

	var res channel_pkg.ChannelListResponse
	res.Channels = channel_pkg.ToChannels(channels)

	return res, nil
}

func (h *ChannelHandler) HandleCreateChannel(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ret, err := h.createChannel(r.Context(), string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, ok := ret.(channel_pkg.CreateChannelResponse)
	if !ok {
		http.Error(w, "잘못된 응답 형식입니다.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error().Err(err).Msg("채널 임대 응답 인코딩 중 오류가 발생했습니다.")
	}
}

func (h *ChannelHandler) HandleRenewChannel(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ret, err := h.renewChannel(r.Context(), string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, ok := ret.(channel_pkg.RenewChannelResponse)
	if !ok {
		http.Error(w, "잘못된 응답 형식입니다.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error().Err(err).Msg("채널 갱신 응답 인코딩 중 오류가 발생했습니다.")
	}
}

func (h *ChannelHandler) HandleListChannels(w http.ResponseWriter, r *http.Request) {

	ret, err := h.listChannels(r.Context(), "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, ok := ret.(channel_pkg.ChannelListResponse)
	if !ok {
		http.Error(w, "잘못된 응답 형식입니다.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error().Err(err).Msg("채널 목록 응답 인코딩 중 오류가 발생했습니다.")
	}
}
