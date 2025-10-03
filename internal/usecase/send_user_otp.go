package usecase

import (
	"context"
	"labra/internal/entity"
)

func (u *UseCase) SendUserOTP(ctx context.Context, login entity.EmailOrPhone) error {
	contact, err := u.contactSvc.GetByValue(ctx, login)
	if err != nil {
		return err
	}

	return u.verifierSvc.SendUserContactVerificationCode(ctx, contact)
}
