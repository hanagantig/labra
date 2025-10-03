package auth

import (
	"context"
	"labra/internal/entity"
	"time"
)

func (s *Service) RefreshSession(ctx context.Context, session entity.Session) (entity.Session, error) {
	rt, err := entity.NewRefreshToken()
	if err != nil {
		return entity.Session{}, err
	}

	// TODO: In the future, we can add a check for the number of active sessions per user and device.
	// If the limit is exceeded, we can revoke the oldest session or all sessions except the current one.

	//TODO: consider DPoP proof key for refresh token rotation
	// https://datatracker.ietf.org/doc/html/draft-ietf-oauth-dpop-08#section-3.4

	// TODO: device id validation not strict if we use DPoP

	newSession := entity.Session{
		AuthTokens: entity.AuthTokens{
			RefreshToken: rt,
		},
		UserUUID:  session.UserUUID,
		SessionID: session.SessionID,
		DeviceID:  session.DeviceID,
		ExpiresAt: time.Now().Add(s.refreshTokenExpiry),
	}

	session = session.Replaced(newSession)

	err = s.tokenRepo.InTransaction(ctx, func(ctx context.Context) error {
		err = s.tokenRepo.UpdateSession(ctx, session)
		if err != nil {
			return err
		}

		err = s.tokenRepo.SaveSession(ctx, newSession)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return entity.Session{}, err
	}

	accessToken, err := entity.NewUserJWT(
		newSession,
		s.accessTokenSecret,
		s.accessTokenExpiry,
	)
	if err != nil {
		return entity.Session{}, err
	}

	newSession.AuthTokens.AccessToken = accessToken

	return newSession, err
}
