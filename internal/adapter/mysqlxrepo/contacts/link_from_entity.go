package contacts

import (
	"context"
	"errors"
)

func (r *Repository) LinkFromEntity(ctx context.Context, fromEntityType, toEntityType string, fromEntityID, toEntityID int) error {
	const query = `INSERT INTO linked_contacts (contact_id, entity_type, entity_id)
						SELECT nl.contact_id, ?, ?
							FROM (
         						SELECT * FROM linked_contacts WHERE entity_id = ? AND entity_type = ?
     						) AS nl;`

	res, err := r.GetConn(ctx).ExecContext(ctx, query, toEntityType, toEntityID, fromEntityID, fromEntityType)
	if err != nil {
		return err
	}

	linkID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	if linkID == 0 {
		return errors.New("failed to link contacts from entity")
	}

	return nil
}
