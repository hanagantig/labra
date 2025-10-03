package unit

import (
	"context"
	"labra/internal/entity"
)

func (r *Repository) GetAll(ctx context.Context) (entity.Units, error) {
	const query = `
		SELECT id, name, unit, description
		FROM units
		ORDER BY name ASC
	`

	var models Units

	conn := r.GetConn(ctx)
	err := conn.SelectContext(ctx, &models, query)
	if err != nil {
		return nil, err
	}

	return models.ToEntity(), nil
}
