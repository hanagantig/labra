package contact

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (s *Service) LinkProfile(ctx context.Context, profileID uuid.UUID, contactID int) error {
	return s.contactRepo.AddLink(ctx, contactID, linkEntityTypeProfile, fmt.Sprint(profileID))
}
