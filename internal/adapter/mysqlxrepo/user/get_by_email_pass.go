package user

import (
	"context"
	"labra/internal/entity"
)

func (r *Repository) GetByCredentials(ctx context.Context, login entity.EmailOrPhone, pass entity.UserPassword) (entity.User, error) {
	const query = `SELECT u.id, u.uuid, u.f_name, u.l_name FROM users as u
						LEFT JOIN linked_contacts as lc ON u.id = lc.entity_id AND lc.entity_type="user"
				  		LEFT JOIN contacts as c ON lc.contact_id = c.id
					WHERE c.value = ? AND u.password = ?;`

	conn := r.GetConn(ctx)

	var model User
	err := conn.GetContext(ctx, &model, query, login.String(), pass)
	if err != nil {
		return entity.User{}, err
	}

	return model.buildEntity(), nil
}
