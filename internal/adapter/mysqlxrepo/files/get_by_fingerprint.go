package files

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) GetByFingerprint(ctx context.Context, fingerprint string) (entity.UploadedFile, error) {
	const query = `SELECT id, file_id, user_id, pipeline_id, fingerprint, file_type, source, status, details, created_at, updated_at 
					FROM uploaded_files WHERE fingerprint = ?`

	conn := r.GetConn(ctx)

	var model UploadedFile
	err := conn.GetContext(ctx, &model, query, fingerprint)
	if err != nil {
		return entity.UploadedFile{}, apperror.ToAppError(err)
	}

	return model.buildEntity(), nil
}
