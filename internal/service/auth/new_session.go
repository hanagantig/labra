package auth

import (
	"context"
	"github.com/google/uuid"
	"labra/internal/entity"
	"time"
)

func (s *Service) NewSession(ctx context.Context, user entity.User, device entity.Device) (entity.Session, error) {
	//TODO: get session by device
	//if existing session is active, use it
	//else create new session

	rToken, err := entity.NewRefreshToken()
	if err != nil {
		return entity.Session{}, err
	}

	session := entity.Session{
		AuthTokens: entity.AuthTokens{
			RefreshToken: rToken,
		},
		UserUUID:  user.Uuid,
		SessionID: uuid.New(),
		ExpiresAt: time.Now().Add(s.refreshTokenExpiry),
	}

	accessToken, err := entity.NewUserJWT(
		session,
		s.accessTokenSecret,
		s.accessTokenExpiry,
	)
	if err != nil {
		return entity.Session{}, err
	}

	err = s.tokenRepo.SaveSession(ctx, session)
	if err != nil {
		return entity.Session{}, err
	}

	session.AuthTokens.AccessToken = accessToken

	return session, nil
}
