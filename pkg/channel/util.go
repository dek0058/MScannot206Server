package channel

import "MScannot206/shared/entity"

func ToChannels(entities []*entity.Channel) []*Channel {
	channels := make([]*Channel, len(entities))
	for i, e := range entities {
		channels[i] = &Channel{
			Id:    e.Id,
			Index: e.Index,
		}
	}
	return channels
}
