package contacts

import (
	"context"
	"github.com/jmoiron/sqlx"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) GetByEntityIDs(ctx context.Context, entityType string, ids []int) (entity.Contacts, error) {
	const query = `SELECT 
    					c.id, 
    					c.type, 
						c.value, 
						c.created_at,
						c.updated_at, 
						c.deleted_at,
						lc.entity_id AS linked_entity_id,
						lc.entity_type AS linked_entity_type
					FROM linked_contacts AS lc
					LEFT JOIN contacts AS c ON c.id = lc.contact_id
					WHERE lc.entity_type = ? AND lc.entity_id IN (?);`

	conn := r.GetConn(ctx)

	if len(ids) == 0 {
		return nil, nil
	}

	q, args, err := sqlx.In(query, entityType, ids)
	if err != nil {
		return nil, err
	}

	var model Contacts
	err = conn.SelectContext(ctx, &model, q, args...)
	if err != nil {
		return nil, apperror.ToAppError(err)
	}

	return model.buildEntity(), nil
}
