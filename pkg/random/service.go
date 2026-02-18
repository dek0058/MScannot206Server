package random

import (
	"context"
	"math/rand/v2"
	"time"
)

func NewRandomService() (*RandomService, error) {
	return &RandomService{}, nil
}

type RandomService struct {
	characterCreateSeed *rand.Rand
}

func (s *RandomService) Start(ctx context.Context) error {
	seed := uint64(time.Now().UnixNano())
	s.characterCreateSeed = rand.New(rand.NewPCG(seed, seed))
	return nil
}

func (s *RandomService) Stop(ctx context.Context) error {
	return nil
}

func (s *RandomService) GetCharacterCreateSeed() *rand.Rand {
	return s.characterCreateSeed
}
