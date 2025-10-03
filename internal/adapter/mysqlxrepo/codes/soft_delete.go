package codes

import (
	"context"
	"labra/internal/entity"
)

func (r *Repository) SoftDelete(ctx context.Context, objType entity.OTPObjectType, objID string) error {
	const query = `UPDATE codes SET deleted_at = now() 
             		WHERE object_type = ? 
             		  AND object_id = ? 
             		  AND used_at is NULL
             		  AND deleted_at IS NULL;`

	_, err := r.GetConn(ctx).ExecContext(ctx, query, objType, objID)
	if err != nil {
		return err
	}

	return nil
}
