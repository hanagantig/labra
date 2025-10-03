package usecase

import (
	"context"
	"labra/internal/entity"
)

func (u *UseCase) VerifyUserContact(ctx context.Context, login entity.EmailOrPhone, otpCode string) (entity.Session, error) {
	contact, err := u.contactSvc.GetByValue(ctx, login)
	if err != nil {
		return entity.Session{}, err
	}

	_, err = u.verifierSvc.VerifyUserContact(ctx, contact, otpCode)
	if err != nil {
		return entity.Session{}, err
	}

	user, err := u.userSvc.GetUserByLogin(ctx, login)
	if err != nil {
		return entity.Session{}, err
	}

	return u.authSvc.NewSession(ctx, user, entity.Device{})
}
