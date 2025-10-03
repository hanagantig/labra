package profile

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) GetByID(ctx context.Context, userID, patientID int) (entity.Profile, error) {
	return s.profileRepo.GetByID(ctx, userID, patientID)
}
