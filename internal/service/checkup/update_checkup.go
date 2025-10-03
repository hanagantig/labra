package checkup

import (
	"context"
	"errors"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (s *Service) UpdateCheckup(ctx context.Context, checkup entity.Checkup, newResults entity.MarkerResults, idsToDelete []int) error {
	err := s.checkupRepo.InTransaction(ctx, func(ctx context.Context) error {
		_, err := s.checkupRepo.UpdateCheckup(ctx, checkup)
		if err != nil && !errors.Is(err, apperror.ErrNotFound) {
			return err
		}

		_, err = s.resultsRepo.Save(ctx, checkup.ID, newResults)
		if err != nil {
			return err
		}

		return s.resultsRepo.DeleteByID(ctx, checkup.ID, idsToDelete)
	})

	return err
}
