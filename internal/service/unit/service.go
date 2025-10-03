package unit

import (
	"context"
	"labra/internal/entity"
)

type Service struct {
	unitRepo unitRepo
}

type unitRepo interface {
	GetNames(ctx context.Context) ([]string, error)
	GetIDByNames(ctx context.Context, names []string) (map[string]int, error)
	GetAll(ctx context.Context) (entity.Units, error)
}

func NewService(ur unitRepo) *Service {
	return &Service{
		unitRepo: ur,
	}
}
