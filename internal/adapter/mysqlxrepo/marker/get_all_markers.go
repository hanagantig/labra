package marker

import (
	"context"
	"labra/internal/entity"
)

func (r *Repository) GetAll(ctx context.Context) (entity.Markers, error) {
	const query = `
		SELECT id, name, ref_range_min, ref_range_max, primary_color
		FROM markers
		ORDER BY name ASC
	`

	var listModel Markers

	conn := r.GetConn(ctx)
	err := conn.SelectContext(ctx, &listModel, query)
	if err != nil {
		return nil, err
	}

	return listModel.BuildEntity(), nil
}
