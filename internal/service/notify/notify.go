package notify

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) Notify(
	ctx context.Context,
	contact entity.Contact,
	templateID string,
	args map[string]string) error {

	channel, err := s.getSender(contact.Type)
	if err != nil {
		return err
	}

	return channel.Send(ctx, contact.Value, templateID, args)
}
