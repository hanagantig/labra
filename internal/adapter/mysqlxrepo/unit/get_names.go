package unit

import (
	"context"
)

func (r *Repository) GetNames(ctx context.Context) ([]string, error) {
	const q = `SELECT DISTINCT name FROM units;`

	conn := r.GetConn(ctx)

	var names []string

	err := conn.SelectContext(ctx, &names, q)
	if err != nil {
		return nil, err
	}

	return names, err
}
