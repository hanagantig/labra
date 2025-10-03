package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (u *UseCase) GetCheckups(ctx context.Context, userUUID, profileUUID uuid.UUID, search string, filter entity.Filter) ([]entity.CheckupResults, error) {
	user, err := u.userSvc.GetByUUID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	if !user.Profiles.HasProfileUUID(profileUUID) {
		return nil, fmt.Errorf("user does not have profile UUID: %s: %w", profileUUID.String(), apperror.ErrNotFound)
	}

	return u.checkupSvc.GetListWithResults(ctx, profileUUID, search, filter)
}
