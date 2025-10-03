package user

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) GetByLogin(ctx context.Context, login entity.EmailOrPhone) (entity.User, error) {
	const query = `SELECT u.id, u.uuid, u.f_name, u.l_name, u.password FROM users as u
						LEFT JOIN linked_contacts as lc ON u.id = lc.entity_id AND lc.entity_type="user"
				  		LEFT JOIN contacts as c ON lc.contact_id = c.id
					WHERE c.value = ?;`

	conn := r.GetConn(ctx)

	var model User
	err := conn.GetContext(ctx, &model, query, login.String())
	if err != nil {
		return entity.User{}, apperror.ToAppError(err)
	}

	return model.buildEntity(), nil
}
