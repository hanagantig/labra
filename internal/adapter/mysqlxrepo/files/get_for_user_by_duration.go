package files

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
	"time"
)

func (r *Repository) GetForUserByDuration(ctx context.Context, userID int, duration time.Duration) ([]entity.UploadedFile, error) {
	const query = `SELECT id, file_id, user_id, pipeline_id, fingerprint, file_type, source, status, details, created_at, updated_at 
					FROM uploaded_files
					WHERE user_id = ? AND created_at >= ?;`

	conn := r.GetConn(ctx)

	var model UploadedFiles
	err := conn.SelectContext(ctx, &model, query, userID, time.Now().Add(-duration))
	if err != nil {
		return nil, apperror.ToAppError(err)
	}

	return model.buildEntity(), nil
}
