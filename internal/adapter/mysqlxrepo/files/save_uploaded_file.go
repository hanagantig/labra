package files

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) SaveUploadedFile(ctx context.Context, file entity.UploadedFile) error {
	const query = `INSERT INTO uploaded_files (user_id, profile_id, file_id, file_type, pipeline_id, fingerprint, status, source, attempts_left) 
					VALUES (:user_id, :profile_id, :file_id, :pipeline_id, :file_type, :fingerprint, :status, :source, 4)`

	model := NewUploadedFileFromEntity(file)
	_, err := r.GetConn(ctx).NamedExecContext(ctx, query, model)
	if err != nil {
		return apperror.ToAppError(err)
	}

	return nil
}
