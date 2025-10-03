package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"labra/internal/entity"
)

func (r *Repository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	const query = `INSERT INTO users(uuid, password, f_name, l_name)
					VALUES (:uuid, :password, :f_name, :l_name)`

	if user.Uuid == uuid.Nil {
		user.Uuid = uuid.New()
	}
	model := NewUserFromEntity(user)

	res, err := r.GetConn(ctx).NamedExecContext(ctx, query, model)
	if err != nil {
		return user, err
	}

	clientID, err := res.LastInsertId()
	if err != nil {
		return user, err
	}

	if clientID == 0 {
		return user, errors.New("failed to insert user")
	}

	model.ID = int(clientID)

	return model.buildEntity(), nil
}
