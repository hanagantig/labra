package results

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"labra/internal/apperror"
)

func (r *Repository) DeleteByID(ctx context.Context, checkupID int, resultIDs []int) error {
	const query = `DELETE FROM checkup_results
		WHERE checkup_id = ? AND id IN (?)`

	conn := r.GetConn(ctx)

	if len(resultIDs) == 0 {
		return nil
	}

	q, args, err := sqlx.In(query, checkupID, resultIDs)
	if err != nil {
		return apperror.ToAppError(err)
	}

	res, err := conn.ExecContext(ctx, q, args...)
	if err != nil {
		return apperror.ToAppError(err)
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return apperror.ToAppError(err)
	}

	if affectedRows <= 0 {
		return fmt.Errorf("no rows found to delete from checkup_results: %w", apperror.ErrNotFound)
	}

	return nil
}
