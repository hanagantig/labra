package lab

import (
	"context"
	"labra/internal/entity"
)

func (r *Repository) GetParserRules(ctx context.Context, labID int) (entity.LabParser, error) {
	return entity.LabParser{}, nil
}
