package checkup

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) GetByID(ctx context.Context, checkupID int) (entity.CheckupResults, error) {
	ch, err := s.checkupRepo.GetByID(ctx, checkupID)
	if err != nil {
		return entity.CheckupResults{}, err
	}

	res, err := s.resultsRepo.GetByCheckup(ctx, checkupID)
	if err != nil {
		return entity.CheckupResults{}, err
	}

	return entity.CheckupResults{
		Checkup: ch,
		Results: res,
	}, nil
}
