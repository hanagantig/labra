package recognizer

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) RunRecognizePipeline(ctx context.Context, f entity.UploadedFile) (string, error) {
	stdID, err := s.dpRepo.Standardize(ctx, f.FileID)
	if err != nil {
		return "", err
	}

	return stdID, nil
}
