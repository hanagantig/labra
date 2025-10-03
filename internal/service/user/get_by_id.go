package user

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) GetByID(ctx context.Context, userID int) (entity.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
