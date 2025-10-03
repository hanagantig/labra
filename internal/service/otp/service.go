package otp

import (
	"context"
	"labra/internal/entity"
	"labra/internal/service"
	"time"
)

type codeRepo interface {
	service.Transactor
	Create(ctx context.Context, code entity.OTPCode) (entity.OTPCode, error)
	UseCode(ctx context.Context, code string, objType entity.OTPObjectType, objID string) error
	SoftDelete(ctx context.Context, objType entity.OTPObjectType, objID string) error
	UnusedCodesCnt(ctx context.Context, objType entity.OTPObjectType, objID string, age time.Duration) (int, error)
	GetByCode(ctx context.Context, code string, objectType entity.OTPObjectType, objectID int) (entity.OTPCode, error)
}

type Service struct {
	codeRepo codeRepo
}

func NewService(cr codeRepo) *Service {
	return &Service{
		codeRepo: cr,
	}
}
