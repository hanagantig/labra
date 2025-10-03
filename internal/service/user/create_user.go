package user

import (
	"context"
	"errors"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (s *Service) CreateUser(ctx context.Context, login entity.EmailOrPhone, user entity.User) (entity.User, error) {
	contact, err := s.contactSvc.GetByValue(ctx, login)
	if err != nil && !errors.Is(err, apperror.ErrNotFound) {
		return entity.User{}, err
	}

	linkedUser := contact.LinkedUser()
	if linkedUser.IsUserType() {
		return entity.User{}, apperror.ErrDuplicateEntity
	}

	err = s.userRepo.InTransaction(ctx, func(ctx context.Context) error {
		var err error

		if contact.ID == 0 {
			contact, err = s.contactSvc.CreateContactByLogin(ctx, login)
			if err != nil {
				return err
			}
		}

		user, err = s.userRepo.Create(ctx, user)
		if err != nil {
			return err
		}

		err = s.contactSvc.LinkUser(ctx, user.ID, contact.ID)
		if err != nil {
			return err
		}

		return err
	})

	if err != nil {
		return user, err
	}

	code, err := s.otpSvc.GenerateOTP(ctx, user.ID, entity.NewUserContactObjectType, contact.IDString())
	if err != nil {
		return user, err
	}

	args := map[string]string{
		"code": code.String(),
	}

	err = s.notifySvc.Notify(ctx, contact, "otp", args)
	if err != nil {
		return user, err
	}

	return user, nil
}
