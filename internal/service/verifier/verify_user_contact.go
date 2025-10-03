package verifier

import (
	"context"
	"errors"
	"fmt"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (s *Service) VerifyUserContact(ctx context.Context, contact entity.Contact, code string) (entity.Contact, error) {
	if len(contact.LinkedEntities) == 0 {
		return entity.Contact{}, fmt.Errorf("no linked entities: %w", apperror.ErrNotFound)
	}

	if contact.VerifiedByUser() {
		return entity.Contact{}, apperror.ErrEntityAlreadyVerified
	}

	otp, err := s.otpSvc.GetOTP(ctx, entity.NewUserContactObjectType, contact.ID, code)
	if err != nil {
		return entity.Contact{}, err
	}

	err = s.otpSvc.VerifyOTP(ctx, otp)
	if err != nil {
		return entity.Contact{}, err
	}

	associatedProfile, err := s.profileSvc.GetAssociatedProfile(ctx, otp.UserID)
	if err != nil && !errors.Is(err, apperror.ErrNotFound) {
		return entity.Contact{}, err
	}

	linkedProfile := contact.LinkedProfile()
	linkedUser := contact.LinkedUser()

	if associatedProfile.ID == 0 {
		//if no profile with associated user and contact has linked profile - migrate profile
		// otherwise - create a profile
		if linkedProfile.ID > 0 {
			err = s.profileSvc.MigrateToUser(ctx, linkedProfile.ID, linkedUser.ID)
			if err != nil {
				return entity.Contact{}, err
			}
			associatedProfile.ID = linkedProfile.ID
		} else {
			p, err := s.profileSvc.CreateForUser(ctx, linkedUser.ID, entity.Profile{
				UserID:        linkedUser.ID,
				CreatorUserID: linkedUser.ID,
			})
			if err != nil {
				return entity.Contact{}, err
			}
			associatedProfile.ID = p.ID
		}
	}

	contact, err = s.contactSvc.VerifyForUser(ctx, otp.UserID, contact, associatedProfile)
	if err != nil {
		return entity.Contact{}, err
	}

	return contact, nil
}
