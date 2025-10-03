package otp

import (
	"context"
	"fmt"
	"labra/internal/entity"
)

func (s *Service) GetOTP(ctx context.Context, objType entity.OTPObjectType, objID int, code string) (entity.OTPCode, error) {
	otp, err := s.codeRepo.GetByCode(ctx, code, objType, objID)
	if err != nil {
		return entity.OTPCode{}, fmt.Errorf("get otp: %w", err)
	}

	return otp, nil
}
