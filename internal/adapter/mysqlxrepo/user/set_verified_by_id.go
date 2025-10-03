package user

import (
	"context"
	"errors"
)

func (r *Repository) SetVerifiedByID(ctx context.Context, userID int) error {
	const query = `UPDATE users SET is_verified = 1 WHERE id = ?`

	res, err := r.GetConn(ctx).ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return errors.New("client does not exist")
	}

	return nil
}
