package results

import (
	"context"
	"errors"
	"labra/internal/entity"
)

func (r *Repository) Save(ctx context.Context, checkupID int, markers entity.MarkerResults) (entity.MarkerResults, error) {
	const query = `INSERT INTO checkup_results 
    	(checkup_id, marker_id, undefined_marker, unit_id, undefined_unit, value, created_at, updated_at) 
		VALUES (:checkup_id, :marker_id, :undefined_marker, :unit_id, :undefined_unit, :value, now(), now());`

	if len(markers) == 0 {
		return nil, nil
	}

	conn := r.GetConn(ctx)

	model := NewResultsFromEntity(checkupID, markers)

	res, err := conn.NamedExecContext(ctx, query, model)
	if err != nil {
		return nil, err
	}

	insertedRows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if insertedRows <= 0 {
		return nil, errors.New("failed to insert checkup_results")
	}

	return markers, nil
}
