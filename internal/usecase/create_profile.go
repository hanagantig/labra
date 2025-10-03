package usecase

import (
	"context"
	"github.com/google/uuid"
	"labra/internal/entity"
)

func (u *UseCase) CreateProfile(ctx context.Context, userUUID uuid.UUID, profile entity.Profile) (entity.Profile, error) {
	user, err := u.userSvc.GetByUUID(ctx, userUUID)
	if err != nil {
		return entity.Profile{}, err
	}

	profile.CreatorUserID = user.ID

	return u.profileSvc.CreateForUser(ctx, user.ID, profile)
}
