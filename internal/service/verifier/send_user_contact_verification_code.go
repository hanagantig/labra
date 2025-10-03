package verifier

import (
	"context"
	"fmt"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (s *Service) SendUserContactVerificationCode(ctx context.Context, contact entity.Contact) error {
	linkedUser := contact.LinkedUser()

	if linkedUser.ID == 0 || !linkedUser.IsUserType() {
		return fmt.Errorf("contact doesn't belong to user: %w", apperror.ErrNotFound)
	}

	code, err := s.otpSvc.GenerateOTP(ctx, linkedUser.ID, entity.NewUserContactObjectType, contact.IDString())
	if err != nil {
		return err
	}

	args := map[string]string{
		"code": code.String(),
	}

	go func() {
		// TODO: log errors
		_ = s.notifySvc.Notify(ctx, contact, "otp", args)
	}()

	return nil
}
