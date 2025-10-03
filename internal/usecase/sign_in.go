package usecase

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (u *UseCase) UserSignIn(ctx context.Context, login entity.EmailOrPhone, pass entity.UserPassword) (entity.Session, error) {
	user, err := u.userSvc.GetUserByLogin(ctx, login)
	if err != nil {
		return entity.Session{}, err
	}

	if !user.Password.IsMatches(pass) {
		return entity.Session{}, apperror.ErrNotFound
	}

	session, err := u.authSvc.NewSession(ctx, user, entity.Device{})
	if err != nil {
		return entity.Session{}, err
	}

	return session, nil
}
