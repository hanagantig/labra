package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (u *UseCase) UploadCheckup(ctx context.Context, userUUID, profileUUID uuid.UUID, fileType string, file []byte) (entity.CheckupResults, error) {
	profiles, err := u.profileSvc.GetUserProfiles(ctx, userUUID)
	if err != nil {
		return entity.CheckupResults{}, fmt.Errorf("unable to get user profiles: %w", err)
	}

	profile := profiles.ProfileByUUID(profileUUID)
	if profile.ID == 0 {
		return entity.CheckupResults{}, fmt.Errorf("profile id not found: %w", apperror.ErrNotFound)
	}

	_, err = u.uploadSvc.UploadFile(ctx, profile, fileType, file)
	if err != nil {
		return entity.CheckupResults{}, fmt.Errorf("unable to upload file: %w", err)
	}

	return entity.CheckupResults{}, nil
}
