package profile

import (
	"context"
	"labra/internal/entity"
)

func (r *Repository) GetByID(ctx context.Context, userID, profileID int) (entity.Profile, error) {
	const query = `SELECT 
    					p.uuid, 
    					p.user_id, 
    					p.creator_user_id, 
    					p.f_name, 
    					p.l_name, 
    					p.gender,
    					up.user_id AS linked_user_id,
    					up.access_level
					FROM user_profiles AS up
					LEFT JOIN profiles AS p ON up.profile_id = p.uuid
					WHERE p.uuid = ? AND p.user_id = ?
					  AND p.deleted_at IS NULL`

	conn := r.GetConn(ctx)
	var model Profile

	err := conn.GetContext(ctx, &model, query, profileID, userID)
	if err != nil {
		return entity.Profile{}, err
	}

	return model.buildEntity(), nil
}
