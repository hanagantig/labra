package otp

import (
	"context"
	"errors"
	"labra/internal/entity"
	"time"
)

const (
	defaultCodeTTL                     = 60 * time.Second
	defaultGenerationFrequencyDuration = time.Minute
)

func (s *Service) GenerateOTP(ctx context.Context, userID int, otpType entity.OTPObjectType, objectID string) (entity.OTPCode, error) {
	code := entity.NewCode(userID, otpType, objectID, defaultCodeTTL)

	cnt, err := s.codeRepo.UnusedCodesCnt(ctx, otpType, objectID, defaultGenerationFrequencyDuration)
	if err != nil {
		return entity.OTPCode{}, err
	}

	if cnt > 0 {
		return entity.OTPCode{}, errors.New("too frequent otp codes")
	}

	var generatedOTP entity.OTPCode

	err = s.codeRepo.InTransaction(ctx, func(ctx context.Context) error {
		err := s.codeRepo.SoftDelete(ctx, otpType, objectID)
		if err != nil {
			return err
		}

		generatedOTP, err = s.codeRepo.Create(ctx, code)

		return err
	})

	if err != nil {
		return entity.OTPCode{}, err
	}

	return generatedOTP, nil
}
