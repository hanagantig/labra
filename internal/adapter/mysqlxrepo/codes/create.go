package codes

import (
	"context"
	"errors"
	"labra/internal/entity"
)

func (r *Repository) Create(ctx context.Context, code entity.OTPCode) (entity.OTPCode, error) {
	const query = `INSERT INTO codes (user_id, object_type, object_id, code, expired_at)
					VALUES (:user_id, :object_type, :object_id, :code, :expired_at)`

	model := NewCodeFromEntity(code)
	res, err := r.GetConn(ctx).NamedExecContext(ctx, query, model)
	if err != nil {
		return code, err
	}

	codeID, err := res.LastInsertId()
	if err != nil {
		return code, err
	}

	if codeID == 0 {
		return code, errors.New("failed to insert code")
	}

	model.ID = uint(codeID)

	return model.buildEntity(), nil
}
