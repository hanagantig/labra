package codes

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) GetByCode(ctx context.Context, code string, objectType entity.OTPObjectType, objectID int) (entity.OTPCode, error) {
	const query = `SELECT id, user_id, object_type, object_id, code, expired_at FROM codes 
					WHERE object_type = ? 
					  AND object_id = ? 
					  AND code = ?
					  AND expired_at > now()
				    ORDER BY created_at
				    LIMIT 1`

	var model Code
	err := r.GetConn(ctx).GetContext(ctx, &model, query, objectType.String(), objectID, code)
	if err != nil {
		return entity.OTPCode{}, apperror.ToAppError(err)
	}

	return model.buildEntity(), nil
}
