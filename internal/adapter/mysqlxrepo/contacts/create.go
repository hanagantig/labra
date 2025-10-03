package contacts

import (
	"context"
	"errors"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) Create(ctx context.Context, contact entity.Contact) (entity.Contact, error) {
	const query = `INSERT INTO contacts (type, value)
					VALUES (:type, :value)`

	model := NewFromEntity(contact)

	res, err := r.GetConn(ctx).NamedExecContext(ctx, query, model)
	if err != nil {
		return entity.Contact{}, apperror.ToAppError(err)
	}

	contactID, err := res.LastInsertId()
	if err != nil {
		return entity.Contact{}, apperror.ToAppError(err)
	}

	if contactID == 0 {
		return entity.Contact{}, apperror.ToAppError(errors.New("failed to insert contact"))
	}

	model.ID = contactID

	return model.buildEntity(), nil
}
