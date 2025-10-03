package checkup

import (
	"context"
	"errors"
	"labra/internal/entity"
)

func (r *Repository) Save(ctx context.Context, e entity.Checkup) (entity.Checkup, error) {
	const query = `INSERT INTO checkups 
    	(title, profile_id, lab_id, date, comment, status, uploaded_file_id, created_at, updated_at) 
		VALUES (:title, :profile_id, :lab_id, :date, :comment, :status, :uploaded_file_id, now(), now());`

	conn := r.GetConn(ctx)

	model := NewCheckupFromEntity(e)

	res, err := conn.NamedExecContext(ctx, query, model)
	if err != nil {
		return entity.Checkup{}, err
	}

	insertedRows, err := res.RowsAffected()
	if err != nil {
		return entity.Checkup{}, err
	}

	if insertedRows <= 0 {
		return entity.Checkup{}, errors.New("failed to insert checkup")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return entity.Checkup{}, err
	}

	e.ID = int(id)

	return e, nil
}
