package file

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) RegisterUploadedFile(ctx context.Context, f entity.UploadedFile) error {
	return s.fileRepo.SaveUploadedFile(ctx, f)
}
