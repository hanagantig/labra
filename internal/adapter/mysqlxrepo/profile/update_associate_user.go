package profile

import (
	"context"
	"fmt"
	"labra/internal/apperror"
)

func (r *Repository) UpdateAssociatedUser(ctx context.Context, profileID, userID int) error {
	const query = `UPDATE profiles SET user_id = ? WHERE id = ?`

	res, err := r.GetConn(ctx).ExecContext(ctx, query, userID, profileID)
	if err != nil {
		return apperror.ToAppError(err)
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return apperror.ToAppError(err)
	}

	if cnt == 0 {
		return fmt.Errorf("can't update associated user: %w", apperror.ErrNotFound)
	}

	return nil
}
