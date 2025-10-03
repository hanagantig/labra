package unit

import (
	"context"
	"labra/internal/entity"
	"labra/pkg/speller"
	"strings"
)

func (s *Service) Identify(ctx context.Context, res entity.CheckupResults) (entity.CheckupResults, error) {
	units, err := s.unitRepo.GetNames(ctx)
	if err != nil {
		return res, err
	}

	spl := speller.New(units, 3)

	correctedMarkers := make([]string, 0, len(res.Results))
	for i := 0; i < len(res.Results); i++ {
		orig := strings.ToLower(res.Results[i].Unit.Name)
		cm := spl.CorrectSpelling(orig)
		if cm != "" {
			correctedMarkers = append(correctedMarkers, cm)
			res.Results[i].Unit.Name = cm
		}
	}

	identifiedMarkers, err := s.unitRepo.GetIDByNames(ctx, correctedMarkers)
	if err != nil {
		return res, err
	}

	for i := 0; i < len(res.Results); i++ {
		if id, ok := identifiedMarkers[res.Results[i].Unit.Name]; ok {
			res.Results[i].Unit.ID = id
		}
	}

	return res, nil
}
