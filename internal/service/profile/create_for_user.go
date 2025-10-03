package profile

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) CreateForUser(ctx context.Context, userID int, profile entity.Profile) (entity.Profile, error) {
	err := s.profileRepo.InTransaction(ctx, func(ctx context.Context) error {
		var err error

		profile, err = s.profileRepo.Upsert(ctx, profile)
		if err != nil {
			return err
		}

		err = s.profileRepo.BindToUser(ctx, profile.ID, userID, "owner")
		if err != nil {
			return err
		}

		// if no associated user - no contacts to link
		if profile.UserID != userID {
			return nil
		}

		err = s.contactSvc.LinkUserContactsToProfile(ctx, profile.ID, userID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.Profile{}, err
	}

	return profile, nil
}
