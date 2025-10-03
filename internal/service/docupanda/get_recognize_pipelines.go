package recognizer

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) GetRecognizePipelineID(ctx context.Context, f entity.UploadedFile) (string, error) {
	if f.FileID == "" {
		return "", nil
	}

	ids, err := s.dpRepo.GetStandardizations(ctx, f.FileID)
	if err != nil {
		return "", err
	}

	if len(ids) > 0 {
		return ids[0], nil
	}

	return "", nil
}
