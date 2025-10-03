package profile

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) GetByAssociatedUserID(ctx context.Context, userID int) (entity.Profile, error) {
	const query = `SELECT 
    					id, 
    					user_id, 
    					creator_user_id, 
    					f_name, 
    					l_name, 
    					gender
					FROM profiles
					WHERE user_id = ?`

	conn := r.GetConn(ctx)
	var model Profile

	err := conn.GetContext(ctx, &model, query, userID)
	if err != nil {
		return entity.Profile{}, apperror.ToAppError(err)
	}

	return model.buildEntity(), nil
}
