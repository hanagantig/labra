package file

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) GetNewFiles(ctx context.Context) ([]entity.UploadedFile, error) {
	return s.fileRepo.GetByStatus(ctx, entity.UploadedFileStatusNew)
}
