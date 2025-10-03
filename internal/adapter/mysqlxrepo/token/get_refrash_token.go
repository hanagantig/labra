package token

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) GetSessionByToken(ctx context.Context, token entity.RefreshToken) (entity.Session, error) {
	const query = `SELECT 
						id, 
						user_uuid,
						session_id,
						token,
						device_id,
						expires_at, 
						created_at, 
						revoked_at, 
						replaced_at,
						replaced_by
					FROM sessions WHERE token = ?`

	model := RefreshToken{}

	conn := r.GetConn(ctx)
	err := conn.GetContext(ctx, &model, query, token.Hash())
	if err != nil {
		return entity.Session{}, apperror.ToAppError(err)
	}

	return model.ToEntity(), nil
}
