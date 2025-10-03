package profile

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) CreateByUser(ctx context.Context, user entity.User, linkedProfile entity.LinkedEntity) (entity.Profile, error) {
	profile := entity.NewProfileFromUser(user)

	if linkedProfile.IsProfileType() {
		profile.ID = linkedProfile.ID
	}

	profile, err := s.profileRepo.Upsert(ctx, profile)
	if err != nil {
		return entity.Profile{}, err
	}

	err = s.profileRepo.UpdateBoundAccess(ctx, profile.ID, "owner", "editor")
	if err != nil {
		return entity.Profile{}, err
	}

	err = s.profileRepo.BindToUser(ctx, profile.ID, user.ID, "owner")
	if err != nil {
		return entity.Profile{}, err
	}

	// if patient not linked with contact - link it
	if !linkedProfile.IsProfileType() {
		err = s.contactSvc.LinkUserContactsToProfile(ctx, profile.ID, user.ID)
		if err != nil {
			return entity.Profile{}, err
		}
	}

	return profile, nil
}
