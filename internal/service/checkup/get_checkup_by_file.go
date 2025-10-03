package checkup

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) GetCheckupByFile(ctx context.Context, f entity.UploadedFile) (entity.Checkup, error) {
	return s.checkupRepo.GetByFile(ctx, f)
}
