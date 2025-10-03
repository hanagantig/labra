package contacts

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) GetLinkedEntities(ctx context.Context, contactID int) ([]entity.LinkedEntity, error) {
	const query = `SELECT 
    					contact_id, entity_id, entity_type, verified_at
					FROM linked_contacts
					WHERE contact_id = ?;`

	conn := r.GetConn(ctx)

	var model LinkedEntities
	err := conn.SelectContext(ctx, &model, query, contactID)
	if err != nil {
		return nil, apperror.ToAppError(err)
	}

	return model.buildEntity(), nil
}
