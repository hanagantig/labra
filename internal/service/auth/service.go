package auth

import (
	"context"
	"labra/internal/entity"
	"labra/internal/service"
	"time"
)

type tokenRepository interface {
	service.Transactor
	SaveSession(ctx context.Context, session entity.Session) error
	UpdateSession(ctx context.Context, session entity.Session) error
	GetSessionByToken(ctx context.Context, token entity.RefreshToken) (entity.Session, error)
}

type Service struct {
	accessTokenSecret  string
	refreshTokenSecret string
	refreshTokenExpiry time.Duration
	accessTokenExpiry  time.Duration
	tokenRepo          tokenRepository
}

func NewService(accessTokenSecret, refreshTokenSecret string, rexp, aexp time.Duration, tr tokenRepository) *Service {
	if rexp.Seconds() <= 0 {
		rexp = time.Hour
	}

	if aexp.Seconds() <= 0 {
		aexp = time.Minute * 15
	}

	return &Service{
		accessTokenSecret:  accessTokenSecret,
		refreshTokenSecret: refreshTokenSecret,
		refreshTokenExpiry: rexp,
		accessTokenExpiry:  aexp,
		tokenRepo:          tr,
	}
}
