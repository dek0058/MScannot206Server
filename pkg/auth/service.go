package auth

import (
	"MScannot206/pkg/auth/session"
	"MScannot206/shared/entity"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/rs/zerolog/log"
)

func NewAuthService() (*AuthService, error) {
	return &AuthService{}, nil
}

type AuthService struct {
	sessionRepo *session.SessionRepository
}

func (s *AuthService) Start(ctx context.Context) error {
	return nil
}

func (s *AuthService) Stop(ctx context.Context) error {
	return nil
}

func (s *AuthService) SetRepositories(
	sessionRepo *session.SessionRepository,
) error {
	var errs error

	s.sessionRepo = sessionRepo
	if sessionRepo == nil {
		errs = errors.Join(errs, session.ErrSessionRepositoryIsNil)
	}

	return errs
}

func (s *AuthService) generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *AuthService) CreateUserSessions(ctx context.Context, user []*entity.User) ([]*entity.UserSession, []*entity.User, error) {
	sessions := make([]*entity.UserSession, 0, len(user))
	failureUsers := make([]*entity.User, 0)

	for _, u := range user {
		token, err := s.generateToken()
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

	err := s.sessionRepo.SaveUserSessions(ctx, sessions)
	if err != nil {
		return nil, nil, err
	}

	return sessions, failureUsers, nil
}

func (s *AuthService) ValidateUserSessions(ctx context.Context, sessions []*entity.UserSession) ([]string, []string, error) {
	return s.sessionRepo.ValidateUserSessions(ctx, sessions)
}
