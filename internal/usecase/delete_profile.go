package usecase

import (
	"context"
)

func (u *UseCase) DeleteProfile(ctx context.Context, userID, profileID int) error {
	return u.profileSvc.SoftDeletePatient(ctx, userID, profileID)
}
