package channel

import "MScannot206/shared/entity"

type CreateChannelResponse struct {
	Channels []*entity.Channel `json:"channels"`
}

type RenewChannelResponse struct {
	Channels []*entity.Channel `json:"channels"`
}

type ChannelListResponse struct {
	Channels []*entity.Channel `json:"channels"`
}
