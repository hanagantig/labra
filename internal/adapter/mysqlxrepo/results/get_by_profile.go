package results

import (
	"context"
	"github.com/jmoiron/sqlx"
	"labra/internal/entity"
)

func (r *Repository) GetByProfile(ctx context.Context, profileID int, filter entity.MarkerFilter, limit int) (entity.MarkerResults, error) {
	const query = `SELECT
    					res.id,
						res.checkup_id,
						res.marker_id,
						m.name AS marker_name,
						m.primary_color,
						res.value,
						res.unit_id,
						u.name AS unit_name,
						res.undefined_unit,
						res.undefined_marker,
						res.created_at
					FROM (SELECT id 
					      FROM checkups WHERE profile_id = ? AND date >= ? AND date <= ?
						  ORDER BY created_at DESC LIMIT ?
					 	) AS ch
					LEFT JOIN checkup_results AS res ON ch.id = res.checkup_id
					LEFT JOIN markers AS m ON res.marker_id = m.id
					LEFT JOIN units AS u ON res.unit_id = u.id
					WHERE res.checkup_id is not null AND m.id IS NOT NULL AND u.id IS NOT NULL;`

	q := query
	args := []any{profileID, filter.From, filter.To, limit}
	var err error

	if len(filter.Names) > 0 {
		q = query + " WHERE m.name IN (?);"
		args = append(args, filter.Names)
		q, args, err = sqlx.In(q, args...)
	}

	conn := r.GetConn(ctx)

	var listModel MarkerResults

	err = conn.SelectContext(ctx, &listModel, q, args...)
	if err != nil {
		return nil, err
	}

	return listModel.BuildEntity(), nil
}
