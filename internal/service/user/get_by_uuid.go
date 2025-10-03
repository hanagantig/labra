package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"labra/internal/entity"
)

func (s *Service) GetByUUID(ctx context.Context, UserUUID uuid.UUID) (entity.UserWithProfiles, error) {
	user, err := s.userRepo.GetByUUID(ctx, UserUUID)
	if err != nil {
		return entity.UserWithProfiles{}, fmt.Errorf("can't get user by uuid: %w", err)
	}

	// user must have at least 1 profile
	profiles, err := s.profileSvc.GetUserProfiles(ctx, UserUUID)
	if err != nil {
		return entity.UserWithProfiles{}, fmt.Errorf("can't get user profiles: %w", err)
	}

	userWithProfiles := entity.UserWithProfiles{
		User:     user,
		Profiles: profiles,
	}

	return userWithProfiles, nil
}
