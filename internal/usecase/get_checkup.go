package usecase

import (
	"context"
	"labra/internal/entity"
)

func (u *UseCase) GetCheckup(ctx context.Context, checkupID int) (entity.CheckupResults, error) {
	return u.checkupSvc.GetByID(ctx, checkupID)
}
