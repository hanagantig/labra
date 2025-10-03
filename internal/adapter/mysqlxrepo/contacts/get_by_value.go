package contacts

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) GetByValue(ctx context.Context, val string) (entity.Contact, error) {
	const query = `SELECT id, type, value, created_at, updated_at, deleted_at 
					FROM contacts
					WHERE value = ?;`

	conn := r.GetConn(ctx)

	var model Contact
	err := conn.GetContext(ctx, &model, query, val)
	if err != nil {
		return entity.Contact{}, apperror.ToAppError(err)
	}

	return model.buildEntity(), nil
}
