package lab

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) GetAll(ctx context.Context) ([]entity.Lab, error) {
	const query = `
		SELECT id, lab_name
		FROM labs
		ORDER BY lab_name ASC
	`

	models := make(Labs, 0)

	conn := r.GetConn(ctx)
	err := conn.SelectContext(ctx, &models, query)
	if err != nil {
		return nil, apperror.ToAppError(err)
	}

	return models.BuildEntity(), nil
}
