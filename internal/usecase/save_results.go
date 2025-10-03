package usecase

import (
	"context"
	"labra/internal/entity"
)

func (u *UseCase) SaveResults(ctx context.Context, checkupResults entity.CheckupResults) error {
	return u.checkupSvc.RegisterCheckup(ctx, checkupResults)
}
