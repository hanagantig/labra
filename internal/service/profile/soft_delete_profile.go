package profile

import (
	"context"
	"fmt"
	"labra/internal/apperror"
)

func (s *Service) SoftDeletePatient(ctx context.Context, userID, patientID int) error {
	patient, err := s.profileRepo.GetByID(ctx, userID, patientID)
	if err != nil {
		return err
	}

	if patient.IsOwnedByUser(userID) {
		return s.profileRepo.SoftDelete(ctx, patientID)
	}

	if !patient.IsBoundToUser(userID) {
		return fmt.Errorf("patient does not belong to user: %w", apperror.ErrNotFound)
	}

	return s.profileRepo.UnBindFromUser(ctx, patientID, userID)
}
