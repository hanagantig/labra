package contact

import (
	"context"
	"fmt"
	"labra/internal/entity"
)

func (s *Service) GetByValue(ctx context.Context, login entity.EmailOrPhone) (entity.Contact, error) {
	contact, err := s.contactRepo.GetByValue(ctx, login.String())
	if err != nil {
		return entity.Contact{}, fmt.Errorf("get contact by value: %w", err)
	}

	linkedEntities, err := s.contactRepo.GetLinkedEntities(ctx, contact.ID)
	if err != nil {
		return entity.Contact{}, fmt.Errorf("get linked entities: %w", err)
	}

	contact.LinkedEntities = linkedEntities

	return contact, err
}
