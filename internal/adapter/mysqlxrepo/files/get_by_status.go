package files

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) GetByStatus(ctx context.Context, status entity.UploadedFileStatus) ([]entity.UploadedFile, error) {
	const query = `SELECT id, file_id, user_id, profile_id, pipeline_id, fingerprint, file_type, source, status, details, created_at, updated_at 
					FROM uploaded_files
					WHERE attempts_left > 0 AND status = ?;`

	conn := r.GetConn(ctx)

	var model UploadedFiles
	err := conn.SelectContext(ctx, &model, query, status)
	if err != nil {
		return nil, apperror.ToAppError(err)
	}

	return model.buildEntity(), nil
}
