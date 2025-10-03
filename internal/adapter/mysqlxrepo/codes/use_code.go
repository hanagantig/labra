package codes

import (
	"context"
	"database/sql"
	"labra/internal/entity"
)

func (r *Repository) UseCode(ctx context.Context, code string, objType entity.OTPObjectType, objID string) error {
	const query = `UPDATE codes SET used_at = now() 
             		WHERE object_type = ? 
             		  AND object_id = ? 
             		  AND code = ?
             		  AND used_at is NULL
             		  AND deleted_at IS NULL;`

	res, err := r.GetConn(ctx).ExecContext(ctx, query, objType, objID, code)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
