package checkup

import (
	"context"
	"github.com/google/uuid"
	"labra/internal/entity"
)

func (s *Service) GetList(ctx context.Context, profileID uuid.UUID) (entity.Checkups, error) {
	return s.checkupRepo.GetCheckupsByUUID(ctx, profileID)
}
