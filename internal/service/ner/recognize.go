package ner

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) Recognize(ctx context.Context, text string) (entity.CheckupResults, error) {
	return s.recognizer.Completions(ctx, text)
}
