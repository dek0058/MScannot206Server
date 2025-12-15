package auth

import (
	"MScannot206/shared/entity"
	"MScannot206/shared/service"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/rs/zerolog/log"
)

func NewAuthService(
	host service.ServiceHost,
) (*AuthService, error) {
	if host == nil {
		return nil, errors.New("host is null")
	}

	return &AuthService{
		host: host,
	}, nil
}

type AuthService struct {
	host service.ServiceHost
}

func (s *AuthService) Init() error {
	return nil
}

func (s *AuthService) Start() error {
	return nil
}

func (s *AuthService) Stop() error {
	return nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *AuthService) CreateUserSessions(user []*entity.User) ([]*entity.UserSession, []*entity.User, error) {
	sessions := make([]*entity.UserSession, 0, len(user))
	failureUsers := make([]*entity.User, 0)

	for _, u := range user {
		token, err := generateToken()
		if err != nil {
			log.Warn().Err(err)
			continue
		}

		session := &entity.UserSession{
			Uid:   u.Uid,
			Token: token,
		}

		sessions = append(sessions, session)
	}

	// TODO: db에 세션 저장 로직 추가 필요

	return sessions, failureUsers, nil
}
