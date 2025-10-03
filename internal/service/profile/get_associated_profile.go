package profile

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) GetAssociatedProfile(ctx context.Context, userID int) (entity.Profile, error) {
	return s.profileRepo.GetByAssociatedUserID(ctx, userID)
}
