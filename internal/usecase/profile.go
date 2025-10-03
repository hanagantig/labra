package usecase

import (
	"context"
	"github.com/google/uuid"
	"labra/internal/entity"
)

func (u *UseCase) Profile(ctx context.Context, profileID uuid.UUID) (entity.MarkerResults, error) {
	return nil, nil
	//return u.checkupSvc.GetUserMarkerResults(ctx, patientID)
}
