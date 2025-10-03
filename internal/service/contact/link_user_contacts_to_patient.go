package contact

import (
	"context"
)

func (s *Service) LinkUserContactsToProfile(ctx context.Context, profileID, userID int) error {
	return s.contactRepo.LinkFromEntity(ctx, linkEntityTypeUser, linkEntityTypeProfile, userID, profileID)
}
