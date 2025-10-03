package checkup

import (
	"context"
	"labra/internal/entity"
)

func (r *Repository) GetByID(ctx context.Context, checkupID int) (entity.Checkup, error) {
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
						LEFT JOIN labs AS l ON c.lab_id = l.id
					    LEFT JOIN profiles AS p ON c.profile_id = p.id
						WHERE c.id = ? ORDER BY c.date DESC`

	conn := r.GetConn(ctx)

	var model Checkup

	err := conn.GetContext(ctx, &model, query, checkupID)
	if err != nil {
		return entity.Checkup{}, err
	}

	return model.BuildEntity(), err
}
