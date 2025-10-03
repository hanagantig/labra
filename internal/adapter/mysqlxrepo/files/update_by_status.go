package files

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) UpdateByStatus(ctx context.Context, file entity.UploadedFile, status entity.UploadedFileStatus) error {
	const query = `UPDATE uploaded_files 
					SET status=?, pipeline_id=?, details=?, attempts_left=attempts_left-1 
					WHERE status=? AND fingerprint=?`

	res, err := r.GetConn(ctx).ExecContext(ctx, query, file.Status, file.PipelineID, file.Details, status, file.Fingerprint)
	if err != nil {
		return apperror.ToAppError(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return apperror.ToAppError(err)
	}

	if rows == 0 {
		return apperror.ErrNotFound
	}

	return nil
}
