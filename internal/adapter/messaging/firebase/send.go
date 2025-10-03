package firebase

import (
	"context"
	"labra/internal/entity"
)

func (r *Repository) Send(ctx context.Context, msg entity.Message) error {
	model := NewMessageFromEntity(msg)
	_, err := r.client.Send(ctx, &model)

	return err
}
