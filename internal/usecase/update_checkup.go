package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (u *UseCase) UpdateCheckup(ctx context.Context, userUUID uuid.UUID, checkupData entity.Checkup, resToAdd entity.MarkerResults, toDelete []int) error {
	profiles, err := u.profileSvc.GetUserProfiles(ctx, userUUID)
	if err != nil {
		return err
	}

	checkupProfile := profiles.ProfileByUUID(checkupData.Profile.Uuid)
	if checkupProfile.ID == 0 {
		return fmt.Errorf("can't find profile for user: %w", apperror.ErrNotFound)
	}

	checkupData.Profile = checkupProfile

	return u.checkupSvc.UpdateCheckup(ctx, checkupData, resToAdd, toDelete)
}
