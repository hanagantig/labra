package identifier

import (
	"context"
	"labra/internal/entity"
)

type Service struct {
	identifiers []identifier
}

type identifier interface {
	Identify(ctx context.Context, res entity.CheckupResults) (entity.CheckupResults, error)
}

func NewService(markers identifier, units identifier, lab identifier, patient identifier) *Service {
	return &Service{
		identifiers: []identifier{
			markers, units, lab, patient,
		},
	}
}
