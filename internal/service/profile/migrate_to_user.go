package profile

import (
	"context"
)

func (s *Service) MigrateToUser(ctx context.Context, profileID, userID int) error {
	err := s.profileRepo.UpdateAssociatedUser(ctx, profileID, userID)
	if err != nil {
		return err
	}

	err = s.profileRepo.UpdateBoundAccess(ctx, profileID, "owner", "editor")
	if err != nil {
		return err
	}

	err = s.profileRepo.BindToUser(ctx, profileID, userID, "owner")
	if err != nil {
		return err
	}

	return nil
}
