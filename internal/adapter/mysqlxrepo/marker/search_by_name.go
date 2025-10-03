package marker

import (
	"context"
	"labra/internal/entity"
)

func (r *Repository) SearchByName(ctx context.Context, search string) (entity.Markers, error) {
	const q = `SELECT id, name FROM markers WHERE name LIKE (?);`

	conn := r.GetConn(ctx)

	var listModel Markers

	err := conn.SelectContext(ctx, &listModel, q, "%"+search+"%")
	if err != nil {
		return nil, err
	}

	return listModel.BuildEntity(), err
}
