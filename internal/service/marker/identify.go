package marker

import (
	"context"
	"labra/internal/entity"
	"labra/pkg/speller"
)

func (s *Service) Identify(ctx context.Context, res entity.CheckupResults) (entity.CheckupResults, error) {
	markers, err := s.markerRepo.GetNames(ctx)
	if err != nil {
		return res, err
	}

	spl := speller.New(markers, 3)

	correctedMarkers := make([]string, 0, len(res.Results))
	for i := 0; i < len(res.Results); i++ {
		orig := res.Results[i].Marker.Name
		cm := spl.CorrectSpelling(orig)
		if cm != "" {
			correctedMarkers = append(correctedMarkers, cm)
			res.Results[i].Marker.Name = cm
		}
	}

	identifiedMarkers, err := s.markerRepo.GetIDByNames(ctx, correctedMarkers)
	if err != nil {
		return res, err
	}

	for i := 0; i < len(res.Results); i++ {
		if id, ok := identifiedMarkers[res.Results[i].Marker.Name]; ok {
			res.Results[i].ID = id
		}
	}

	return res, nil
}
