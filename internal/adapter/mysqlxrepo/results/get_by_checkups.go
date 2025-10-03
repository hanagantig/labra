package results

import (
	"context"
	"github.com/jmoiron/sqlx"
	"labra/internal/entity"
)

func (r *Repository) GetByCheckups(ctx context.Context, checkupIDs, makerIDs []int) (map[int]entity.MarkerResults, error) {
	const q = `SELECT
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
					WHERE ch.id IN (?) AND res.id IS NOT NULL`

	query, args, err := sqlx.In(q, checkupIDs)
	if err != nil {
		return nil, err
	}

	if len(makerIDs) != 0 {
		mq, margs, err := sqlx.In("AND m.id IN (?)", makerIDs)
		if err != nil {
			return nil, err
		}

		args = append(args, margs...)
		query += " " + mq
	}

	conn := r.GetConn(ctx)

	var listModel MarkerResults

	err = conn.SelectContext(ctx, &listModel, query, args...)
	if err != nil {
		return nil, err
	}

	groupedByCheckup := make(map[int]MarkerResults)
	for _, res := range listModel {
		groupedByCheckup[res.CheckupID] = append(groupedByCheckup[res.CheckupID], res)
	}

	results := make(map[int]entity.MarkerResults)
	for checkupID, res := range groupedByCheckup {
		if len(res) == 0 {
			continue
		}

		results[checkupID] = res.BuildEntity()
	}

	return results, nil
}
