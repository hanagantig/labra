package contact

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) GetForProfiles(ctx context.Context, profileIDs []int) (entity.Contacts, error) {
	return s.contactRepo.GetByEntityIDs(ctx, linkEntityTypeProfile, profileIDs)
}
