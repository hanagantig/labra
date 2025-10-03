package checkup

import (
	"context"
	"labra/internal/entity"
	"time"
)

const defaultCheckupLimit = 3

func (s *Service) GetUserMarkerResults(ctx context.Context, profileID int, from, to time.Time, names []string) (entity.MarkerResults, error) {
	mFilter := entity.MarkerFilter{
		From:  from,
		To:    to,
		Names: names,
	}

	return s.resultsRepo.GetByProfile(ctx, profileID, mFilter, defaultCheckupLimit)
}
