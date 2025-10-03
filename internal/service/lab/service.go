package lab

import (
	"context"
	"labra/internal/entity"
)

type labRepo interface {
	GetAll(ctx context.Context) ([]entity.Lab, error)
}

type Service struct {
	labRepo labRepo
}

func NewService(lb labRepo) *Service {
	return &Service{
		labRepo: lb,
	}
}
