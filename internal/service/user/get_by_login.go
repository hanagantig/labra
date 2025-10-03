package user

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) GetUserByLogin(ctx context.Context, login entity.EmailOrPhone) (entity.User, error) {
	return s.userRepo.GetByLogin(ctx, login)
}
