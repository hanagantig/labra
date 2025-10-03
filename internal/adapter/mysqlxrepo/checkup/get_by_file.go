package checkup

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) GetByFile(ctx context.Context, f entity.UploadedFile) (entity.Checkup, error) {
	const query = `SELECT 
    				c.id AS id, 
					uploaded_file_id, 
					profile_id, 
					l.lab_name, 
					lab_id, 
					p.f_name, 
					p.l_name, 
    				date 
				   FROM checkups AS c 
						LEFT JOIN labs AS l ON c.lab_id = l.id
					    LEFT JOIN profiles AS p ON c.profile_id = p.id
						WHERE c.uploaded_file_id = ? AND c.profile_id = ?`

	conn := r.GetConn(ctx)

	var model Checkup

	err := conn.GetContext(ctx, &model, query, f.ID, f.ProfileID)
	if err != nil {
		return entity.Checkup{}, apperror.ToAppError(err)
	}

	return model.BuildEntity(), nil
}
