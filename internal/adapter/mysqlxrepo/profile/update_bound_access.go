package profile

import (
	"context"
	"labra/internal/apperror"
)

func (r *Repository) UpdateBoundAccess(ctx context.Context, profileID int, fromAccessLevel, toAccessLevel string) error {
	const query = `UPDATE user_profiles SET access_level = ? 
                     WHERE profile_id = ? AND access_level = ?;`

	_, err := r.GetConn(ctx).ExecContext(ctx, query, toAccessLevel, profileID, fromAccessLevel)
	if err != nil {
		return apperror.ToAppError(err)
	}

	return nil
}
