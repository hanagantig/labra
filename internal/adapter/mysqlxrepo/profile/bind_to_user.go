package profile

import (
	"context"
	"errors"
	"labra/internal/apperror"
)

func (r *Repository) BindToUser(ctx context.Context, profileID, userID int, role string) error {
	const query = `INSERT INTO user_profiles (user_id, profile_id, access_level)
					VALUES (?, ?, ?)`

	res, err := r.GetConn(ctx).ExecContext(ctx, query, userID, profileID, role)
	if err != nil {
		return apperror.ToAppError(err)
	}

	relID, err := res.LastInsertId()
	if err != nil {
		return apperror.ToAppError(err)
	}

	if relID == 0 {
		return apperror.ToAppError(errors.New("failed to bind user to patient"))
	}

	return nil
}
