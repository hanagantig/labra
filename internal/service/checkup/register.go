package checkup

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) RegisterCheckup(ctx context.Context, checkup entity.CheckupResults) error {
	err := s.checkupRepo.InTransaction(ctx, func(ctx context.Context) error {
		check, err := s.checkupRepo.Save(ctx, checkup.Checkup)
		if err != nil {
			return err
		}

		_, err = s.resultsRepo.Save(ctx, check.ID, checkup.Results)

		return err
	})

	return err
}
