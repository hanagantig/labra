package profile

import (
	"context"
	"labra/internal/entity"
)

func (r *Repository) GetByUserID(ctx context.Context, userID int) (entity.Profiles, error) {
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
					WHERE up.user_id = ? 
					  AND p.id IS NOT NULL
					  AND p.deleted_at IS NULL`

	conn := r.GetConn(ctx)
	var listModel Profiles

	err := conn.SelectContext(ctx, &listModel, query, userID)
	if err != nil {
		return nil, err
	}

	return listModel.buildEntity(), nil
}
