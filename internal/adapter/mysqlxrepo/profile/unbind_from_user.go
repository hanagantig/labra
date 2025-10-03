package profile

import (
	"context"
	"labra/internal/apperror"
)

func (r *Repository) UnBindFromUser(ctx context.Context, patientID, userID int) error {
	const query = `DELETE FROM user_patients WHERE user_id = ? AND patient_id = ?`

	_, err := r.GetConn(ctx).ExecContext(ctx, query, userID, patientID)
	if err != nil {
		return apperror.ToAppError(err)
	}

	return nil
}
