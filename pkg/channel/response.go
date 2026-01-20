package channel

import "MScannot206/shared/entity"

type FetchChannelsResponse struct {
	Channels []*entity.Channel `json:"channels"`
}
