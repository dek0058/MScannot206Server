package channel

import (
	"MScannot206/shared/entity"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const LeaseDuration = 5 * time.Minute

func NewChannelService() (*ChannelService, error) {
	return &ChannelService{}, nil
}

type ChannelService struct {
	channelRepo *ChannelMongoRepository
}

func (s *ChannelService) Start(ctx context.Context) error {
	s.startCleanup(ctx, 1*time.Minute)
	return nil
}

func (s *ChannelService) Stop(ctx context.Context) error {
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

func (s *ChannelService) Acquire(ctx context.Context, channelId string) (*entity.Channel, error) {
	nextIndex, err := s.channelRepo.GetNextSequence(ctx)
	if err != nil {
		return nil, err
	}

	newChannel := &entity.Channel{
		ID:        channelId,
		Index:     nextIndex,
		LeaseID:   uuid.NewString(),
		ExpiresAt: time.Now().Add(LeaseDuration),
	}

	if err := s.channelRepo.CreateLease(ctx, *newChannel); err != nil {
		return nil, err
	}

	return newChannel, nil
}

func (s *ChannelService) Renew(ctx context.Context, leaseID string) (*entity.Channel, error) {
	newExpiry := time.Now().UTC().Add(LeaseDuration)
	renewedChannel, err := s.channelRepo.RenewLease(ctx, leaseID, newExpiry)
	if err != nil {
		return nil, err
	}

	return renewedChannel, nil
}

func (s *ChannelService) GetChannels(ctx context.Context) ([]*entity.Channel, error) {
	return s.channelRepo.GetAllActiveLeases(ctx)
}

func (s *ChannelService) startCleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			<-ticker.C
			log.Info().Msg("만료된 채널 정리 작업 시작합니다.")

			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			s.runCleanup(ctx)
			cancel()
		}
	}()
}

func (s *ChannelService) runCleanup(ctx context.Context) {
	expired, err := s.channelRepo.FindExpiredLeases(ctx, time.Now())
	if err != nil {
		log.Error().Err(err).Msg("만료된 채널 조회 중 오류가 발생했습니다.")
		return
	}

	if len(expired) == 0 {
		return
	}

	var expiredIds []string
	for _, ch := range expired {
		expiredIds = append(expiredIds, ch.LeaseID)
		// TODO: 재활용 로직 추가
	}

	deletedCount, err := s.channelRepo.DeleteLeases(ctx, expiredIds)
	if err != nil {
		log.Error().Err(err).Msg("만료된 채널 삭제 중 오류가 발생했습니다.")
		return
	}

	log.Info().Int64("count", deletedCount).Msg("만료된 채널을 정리했습니다.")
}
