package ner

import (
	"context"
	"labra/internal/entity"
)

type recognizer interface {
	Completions(ctx context.Context, message string) (entity.CheckupResults, error)
}

type Service struct {
	recognizer recognizer
}

func NewService(recognizer recognizer) *Service {
	return &Service{recognizer: recognizer}
}
