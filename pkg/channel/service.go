package channel

import (
	"MScannot206/shared/entity"
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

const LeaseDuration = 30 * time.Minute

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

func (s *ChannelService) Create(ctx context.Context, channelId string) (*entity.Channel, error) {
	existingChannel, err := s.channelRepo.FindChannelByID(ctx, channelId)
	if err != nil {
		return nil, err
	}

	if existingChannel != nil {
		log.Info().Str("channel_id", channelId).Msg("이미 존재하는 채널입니다. 갱신을 수행합니다.")
		return s.Renew(ctx, channelId)
	}

	nextIndex, err := s.channelRepo.PopRecyclableIndex(ctx)
	if err != nil {
		return nil, err
	}

	if nextIndex == 0 {
		nextIndex, err = s.channelRepo.GetNextSequence(ctx)
		if err != nil {
			return nil, err
		}
	}

	newChannel := &entity.Channel{
		Id:        channelId,
		Index:     nextIndex,
		ExpiresAt: time.Now().Add(LeaseDuration),
	}

	if err := s.channelRepo.CreateChannel(ctx, *newChannel); err != nil {
		if nextIndex != 0 {
			// 재활용한 인덱스를 다시 넣어줌
			if pushErr := s.channelRepo.PushRecyclableIndex(ctx, nextIndex); pushErr != nil {
				log.Error().Err(pushErr).Int("index", nextIndex).Msg("채널 인덱스를 재활용 목록에 다시 추가하는데 실패했습니다.")
			}
		}
		return nil, err
	}

	return newChannel, nil
}

func (s *ChannelService) Renew(ctx context.Context, channelId string) (*entity.Channel, error) {
	newExpiry := time.Now().UTC().Add(LeaseDuration)
	renewedChannel, err := s.channelRepo.RenewChannel(ctx, channelId, newExpiry)
	if err != nil {
		return nil, err
	}

	return renewedChannel, nil
}

func (s *ChannelService) GetChannels(ctx context.Context) ([]*entity.Channel, error) {
	return s.channelRepo.GetAllActiveChannels(ctx)
}

func (s *ChannelService) startCleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			<-ticker.C
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			s.runCleanup(ctx)
			cancel()
		}
	}()
}

func (s *ChannelService) runCleanup(ctx context.Context) {
	expired, err := s.channelRepo.FindExpiredChannels(ctx, time.Now())
	if err != nil {
		log.Error().Err(err).Msg("만료된 채널 조회 중 오류가 발생했습니다.")
		return
	}

	if len(expired) == 0 {
		return
	}

	log.Info().Msg("만료된 채널 정리 작업 시작합니다.")

	var expiredIds []string
	for _, ch := range expired {
		expiredIds = append(expiredIds, ch.Id)

		if err := s.channelRepo.PushRecyclableIndex(ctx, ch.Index); err != nil {
			log.Error().Err(err).Int("index", ch.Index).Msg("채널 인덱스 재활용 목록에 추가 중 오류 발생했습니다.")
		} else {
			log.Info().Int("index", ch.Index).Msg("채널 인덱스를 재활용 목록에 추가했습니다.")
		}
	}

	deletedCount, err := s.channelRepo.DeleteChannels(ctx, expiredIds)
	if err != nil {
		log.Error().Err(err).Msg("만료된 채널 삭제 중 오류가 발생했습니다.")
		return
	}

	log.Info().Int64("count", deletedCount).Msg("만료된 채널을 정리했습니다.")
}
