package token

import (
	"context"
	"errors"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) UpdateSession(ctx context.Context, session entity.Session) error {
	const query = `UPDATE sessions SET
						expires_at = :expires_at,
						replaced_by = :replaced_by,
						replaced_at = :replaced_at,
						revoked_at = :revoked_at
					WHERE session_id = :session_id 
					  AND token = :token`

	model := NewRefreshTokenFromEntity(session)

	conn := r.GetConn(ctx)
	res, err := conn.NamedExecContext(ctx, query, model)
	if err != nil {
		return apperror.ToAppError(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return apperror.ToAppError(err)
	}

	if affected <= 0 {
		return apperror.ToAppError(errors.New("failed to update session"))
	}

	return nil
}
