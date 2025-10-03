package parser

import (
	"context"
	"labra/internal/entity"
)

type loader interface {
	GetParserRules(ctx context.Context, labID int) (entity.LabParser, error)
}

type Service struct {
	loader  loader
	parsers map[int]entity.LabParser
}

func NewService(l loader) *Service {
	return &Service{
		loader: l,
	}
}
