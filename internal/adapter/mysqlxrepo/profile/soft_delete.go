package profile

import (
	"context"
	"labra/internal/apperror"
)

func (r *Repository) SoftDelete(ctx context.Context, patientID int) error {
	q := `UPDATE patients SET deleted_at = NOW() 
                WHERE id = ?;`

	conn := r.GetConn(ctx)
	_, err := conn.ExecContext(ctx, q, patientID)
	if err != nil {
		return apperror.ToAppError(err)
	}

	return nil
}
