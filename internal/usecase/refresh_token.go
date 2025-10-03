package usecase

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (u *UseCase) RefreshUserToken(ctx context.Context, refreshToken entity.RefreshToken) (entity.Session, error) {
	existingSession, err := u.authSvc.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return entity.Session{}, err
	}

	if !existingSession.IsActive() {
		return entity.Session{}, apperror.ErrUnauthorized
	}

	_, err = u.userSvc.GetByUUID(ctx, existingSession.UserUUID)
	if err != nil {
		return entity.Session{}, err
	}

	refreshedSession, err := u.authSvc.RefreshSession(ctx, existingSession)
	if err != nil {
		return entity.Session{}, err
	}

	return refreshedSession, nil
}
