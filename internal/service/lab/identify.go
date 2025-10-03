package lab

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) Identify(ctx context.Context, res entity.CheckupResults) (entity.CheckupResults, error) {
	return res, nil
}
