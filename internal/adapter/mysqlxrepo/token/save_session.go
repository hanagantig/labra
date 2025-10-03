package token

import (
	"context"
	"errors"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) SaveSession(ctx context.Context, session entity.Session) error {
	const query = `INSERT INTO sessions(user_uuid, session_id, token, device_id, expires_at, replaced_by, replaced_at, revoked_at, created_at)
					VALUES (:user_uuid, :session_id, :token, :device_id, :expires_at, :replaced_by, :replaced_at, :revoked_at, :created_at)`

	model := NewRefreshTokenFromEntity(session)

	conn := r.GetConn(ctx)
	res, err := conn.NamedExecContext(ctx, query, model)
	if err != nil {
		return apperror.ToAppError(err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return apperror.ToAppError(err)
	}

	if lastID <= 0 {
		return apperror.ToAppError(errors.New("failed to insert refresh token"))
	}

	return nil
}
