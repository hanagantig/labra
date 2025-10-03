package usecase

import (
	"context"
	"github.com/google/uuid"
	"labra/internal/entity"
)

func (u *UseCase) GetUserProfiles(ctx context.Context, userUUID uuid.UUID) (entity.Profiles, error) {
	return u.profileSvc.GetUserProfiles(ctx, userUUID)
}
