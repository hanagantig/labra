package checkup

import (
	"context"
	"github.com/google/uuid"
	"labra/internal/entity"
)

func (s *Service) GetListWithResults(
	ctx context.Context,
	profileID uuid.UUID,
	search string,
	filter entity.Filter,
) ([]entity.CheckupResults, error) {

	var foundMarkers entity.Markers
	var err error

	if search != "" {
		foundMarkers, err = s.mr.SearchByName(ctx, search)
		if err != nil {
			return nil, err
		}

		if len(foundMarkers) == 0 {
			return nil, nil
		}
	}

	checkups, err := s.checkupRepo.GetCheckupsByUUID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	if len(checkups) == 0 {
		return nil, nil
	}

	checkupResults, err := s.resultsRepo.GetByCheckups(ctx, checkups.IDs(), foundMarkers.IDs())
	if err != nil {
		return nil, err
	}

	results := make([]entity.CheckupResults, 0, len(foundMarkers))
	for _, checkup := range checkups {
		if len(checkupResults[checkup.ID]) == 0 {
			continue
		}

		results = append(results, entity.CheckupResults{
			Checkup: checkup,
			Results: checkupResults[checkup.ID],
		})
	}

	return results, nil
}
