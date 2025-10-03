package contacts

import (
	"context"
	"errors"
)

func (r *Repository) AddLink(ctx context.Context, contactID int, entityType string, entityID string) error {
	const query = `INSERT INTO linked_contacts (contact_id, entity_type, entity_id)
					VALUES (?, ?, ?)`

	res, err := r.GetConn(ctx).ExecContext(ctx, query, contactID, entityType, entityID)
	if err != nil {
		return err
	}

	linkID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	if linkID == 0 {
		return errors.New("failed to insert linked contact")
	}

	return nil
}
