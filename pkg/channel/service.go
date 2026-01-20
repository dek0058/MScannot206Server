package channel

import (
	"MScannot206/shared/entity"
	"MScannot206/shared/service"
	"context"
	"errors"
)

func NewChannelService() (*ChannelService, error) {
	return &ChannelService{}, nil
}

type ChannelService struct {
	host service.ServiceHost

	channelRepo *ChannelMongoRepository
}

func (s *ChannelService) Start() error {
	return nil
}

func (s *ChannelService) Stop() error {
	return nil
}

func (s *ChannelService) SetRepositories(
	channelRepo *ChannelMongoRepository,
) error {
	var errs error

	s.channelRepo = channelRepo
	if channelRepo == nil {
		errs = errors.Join(errs, ErrChannelRepositoryIsNil)
	}

	return errs
}

func (s *ChannelService) AddChannel(ctx context.Context, id string) error {
	return s.channelRepo.AddChannel(ctx, id)
}

func (s *ChannelService) RemoveChannel(ctx context.Context, id string) error {
	return s.channelRepo.RemoveChannel(ctx, id)
}

func (s *ChannelService) GetChannels(ctx context.Context) ([]*entity.Channel, error) {
	return s.channelRepo.GetChannels(ctx)
}
