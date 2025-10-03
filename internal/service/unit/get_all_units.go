package unit

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) GetAllUnits(ctx context.Context) (entity.Units, error) {
	return s.unitRepo.GetAll(ctx)
}
