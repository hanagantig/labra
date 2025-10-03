package contacts

import (
	"context"
	"errors"
	"labra/internal/entity"
)

func (r *Repository) SetVerifiedByEntity(ctx context.Context, linkedEntity entity.LinkedEntity) error {
	const query = `UPDATE linked_contacts SET verified_at=now() 
                       WHERE contact_id = ?
                       AND entity_type = ?
                       AND entity_id = ?;`

	res, err := r.GetConn(ctx).
		ExecContext(
			ctx, query,
			linkedEntity.ContactID,
			linkedEntity.Type,
			linkedEntity.ID,
		)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return errors.New("contact entity does not exist")
	}

	return nil
}
