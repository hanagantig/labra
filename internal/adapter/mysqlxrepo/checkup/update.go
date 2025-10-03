package checkup

import (
	"context"
	"fmt"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) UpdateCheckup(ctx context.Context, e entity.Checkup) (entity.Checkup, error) {
	const query = `UPDATE checkups
		SET title = :title,
			profile_id = :profile_id,
			lab_id = :lab_id,
			date = :date,
			comment = :comment,
			status = :status
		WHERE id = :id;`

	conn := r.GetConn(ctx)

	model := NewCheckupFromEntity(e)

	res, err := conn.NamedExecContext(ctx, query, model)
	if err != nil {
		return entity.Checkup{}, err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return entity.Checkup{}, err
	}

	if affectedRows <= 0 {
		return entity.Checkup{}, fmt.Errorf("no rows updated for checkup: %w", apperror.ErrNotFound)
	}

	return e, nil
}
