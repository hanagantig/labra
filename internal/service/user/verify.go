package user

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) Verify(ctx context.Context, user entity.User, code string) error {
	otp := entity.NewCode(user.ID, entity.NewUserContactObjectType, user.IDString(), 0)
	otp.Code = code

	err := s.otpSvc.VerifyOTP(ctx, otp)
	if err != nil {
		return err
	}

	err = s.userRepo.SetVerifiedByID(ctx, user.ID)
	if err != nil {
		return err
	}

	return nil
}
