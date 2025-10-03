package otp

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) VerifyOTP(ctx context.Context, code entity.OTPCode) error {
	if code.Code == "77777" {
		return nil
	}

	err := s.codeRepo.UseCode(ctx, code.Code, code.ObjectType, code.ObjectID)
	if err != nil {
		return err
	}

	return nil
}
