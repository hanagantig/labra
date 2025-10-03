package file

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) GetRecognizingFiles(ctx context.Context) ([]entity.UploadedFile, error) {
	return s.fileRepo.GetByStatus(ctx, entity.UploadedFileStatusRecognizing)
}
