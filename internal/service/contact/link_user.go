package contact

import (
	"context"
	"fmt"
)

func (s *Service) LinkUser(ctx context.Context, userID int, contactID int) error {
	return s.contactRepo.AddLink(ctx, contactID, linkEntityTypeUser, fmt.Sprint(userID))
}
