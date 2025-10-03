package auth

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) GetSessionByToken(ctx context.Context, refreshToken entity.RefreshToken) (entity.Session, error) {
	return s.tokenRepo.GetSessionByToken(ctx, refreshToken)
}
