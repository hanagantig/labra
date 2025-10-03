package unit

import (
	"context"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) GetIDByNames(ctx context.Context, names []string) (map[string]int, error) {
	const q = `SELECT id, name FROM units WHERE name IN (?);`

	conn := r.GetConn(ctx)

	query, args, err := sqlx.In(q, names)
	if err != nil {
		return nil, err
	}

	res := []struct {
		ID   int
		Name string
	}{}

	err = conn.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}

	mappedNames := make(map[string]int, len(res))
	for _, r := range res {
		mappedNames[r.Name] = r.ID
	}

	return mappedNames, err
}
