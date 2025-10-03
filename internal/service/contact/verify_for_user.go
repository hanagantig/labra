package contact

import (
	"context"
	"fmt"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (s *Service) VerifyForUser(ctx context.Context, userID int, contact entity.Contact, associatedProfile entity.Profile) (entity.Contact, error) {
	linkedUser := contact.LinkedUser()

	if userID != linkedUser.ID {
		return entity.Contact{}, fmt.Errorf("invalid user id: %w", apperror.ErrNotFound)
	}

	err := s.contactRepo.SetVerifiedByEntity(ctx, linkedUser)
	if err != nil {
		return contact, err
	}

	if associatedProfile.ID > 0 {
		linkedProfile := entity.LinkedEntity{
			ID:        associatedProfile.ID,
			ContactID: contact.ID,
			Type:      entity.LinkedContactEntityProfileType,
		}

		err = s.contactRepo.SetVerifiedByEntity(ctx, linkedProfile)
		if err != nil {
			return contact, err
		}
	}

	return contact, nil
}
