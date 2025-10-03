package marker

import (
	"context"
	"labra/internal/entity"
)

type markerRepo interface {
	GetNames(ctx context.Context) ([]string, error)
	GetIDByNames(ctx context.Context, names []string) (map[string]int, error)
	GetAll(ctx context.Context) (entity.Markers, error)
}

type Service struct {
	markerRepo markerRepo
}

func NewService(mr markerRepo) *Service {
	return &Service{
		markerRepo: mr,
	}
}

func (s *Service) GetAllMarkers(ctx context.Context) ([]entity.Marker, error) {
	return s.markerRepo.GetAll(ctx)
}
