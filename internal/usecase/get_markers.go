package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"labra/internal/entity"
	"time"
)

func (u *UseCase) GetMarkers(ctx context.Context, userUUID, profileUUID uuid.UUID, from, to time.Time, names []string) (entity.MarkerResults, error) {
	user, err := u.userSvc.GetByUUID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	profile := user.Profiles.ProfileByUUID(profileUUID)
	if profile.ID == 0 {
		return nil, fmt.Errorf("user does not have profile UUID: %s", profileUUID.String())
	}

	return u.checkupSvc.GetUserMarkerResults(ctx, profile.ID, from, to, names)
}
