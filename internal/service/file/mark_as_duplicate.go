package file

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) MarkAsDuplicate(ctx context.Context, f entity.UploadedFile) error {
	initialStatus := f.Status
	f.Status = entity.UploadedFileStatusDuplicated

	return s.fileRepo.UpdateByStatus(ctx, f, initialStatus)
}
