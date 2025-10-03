package user

import (
	"context"
	"github.com/google/uuid"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) GetByUUID(ctx context.Context, userUUID uuid.UUID) (entity.User, error) {
	const query = `SELECT u.id, u.uuid, u.f_name, u.l_name FROM users as u
						LEFT JOIN linked_contacts as lc ON u.id = lc.entity_id AND lc.entity_type="user"
				  		LEFT JOIN contacts as c ON lc.contact_id = c.id
					WHERE u.uuid = ?;`

	conn := r.GetConn(ctx)

	var model User
	err := conn.GetContext(ctx, &model, query, userUUID.String())
	if err != nil {
		return entity.User{}, apperror.ToAppError(err)
	}

	return model.buildEntity(), nil
}
