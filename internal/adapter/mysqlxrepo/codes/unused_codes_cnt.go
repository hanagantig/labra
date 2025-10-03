package codes

import (
	"context"
	"labra/internal/entity"
	"time"
)

func (r *Repository) UnusedCodesCnt(ctx context.Context, objType entity.OTPObjectType, objID string, age time.Duration) (int, error) {
	const query = `SELECT count(*) FROM codes 
                	WHERE object_type = ? 
                	  AND object_id = ? 
                	  AND created_at > ?
                	  AND used_at IS NULL 
                	  AND deleted_at IS NULL;`

	var count int

	createdSince := time.Now().Add(-age)
	err := r.GetConn(ctx).GetContext(ctx, &count, query, objType, objID, createdSince)
	if err != nil {
		return 0, err
	}

	return count, nil
}
