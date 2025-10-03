package checkup

import (
	"context"
	"github.com/google/uuid"
	"labra/internal/entity"
)

func (r *Repository) GetCheckupsByUUID(ctx context.Context, profileID uuid.UUID) (entity.Checkups, error) {
	const query = `SELECT 
    					c.id, 
    					c.title, 
    					p.uuid as profile_uuid, 
    					profile_id, 
    					l.lab_name, 
    					lab_id, 
    					p.f_name, 
    					p.l_name, 
    					c.uploaded_file_id,
    					c.status,
    					date 
					FROM checkups AS c 
					    LEFT JOIN profiles AS p ON c.profile_id = p.id
						LEFT JOIN labs AS l ON c.lab_id = l.id
						WHERE p.uuid = ? ORDER BY c.date DESC
						LIMIT 50`

	conn := r.GetConn(ctx)

	var model Checkups

	err := conn.SelectContext(ctx, &model, query, profileID)
	if err != nil {
		return nil, err
	}

	return model.BuildEntity(), err
}
