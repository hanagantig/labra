package results

import (
	"context"
	"labra/internal/entity"
)

func (r *Repository) GetByCheckup(ctx context.Context, checkupID int) (entity.MarkerResults, error) {
	const query = `SELECT
						res.id,
						res.checkup_id,
						res.marker_id,
						m.name AS marker_name,
						res.value,
						res.unit_id,
						u.name AS unit_name,
						res.undefined_unit,
						res.undefined_marker,
						res.created_at
					FROM checkups AS ch
					LEFT JOIN checkup_results AS res ON ch.id = res.checkup_id
					LEFT JOIN markers AS m ON res.marker_id = m.id
					LEFT JOIN units AS u ON res.unit_id = u.id
					WHERE ch.id=?`

	conn := r.GetConn(ctx)

	var listModel MarkerResults

	err := conn.SelectContext(ctx, &listModel, query, checkupID)
	if err != nil {
		return nil, err
	}

	return listModel.BuildEntity(), nil
}
