package channel

import "MScannot206/shared/entity"

type AcquireChannelResponse struct {
	Channels []*entity.Channel `json:"channels"`
}

type RenewChannelResponse struct {
	Channels []*entity.Channel `json:"channels"`
}
