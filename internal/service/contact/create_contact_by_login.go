package contact

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) CreateContactByLogin(ctx context.Context, login entity.EmailOrPhone) (entity.Contact, error) {
	return s.contactRepo.Create(ctx, entity.NewContactFromLogin(login))
}
