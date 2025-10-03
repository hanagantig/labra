package profile

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (s *Service) GetUserProfiles(ctx context.Context, userUUID uuid.UUID) (entity.Profiles, error) {
	profiles, err := s.profileRepo.GetByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	if len(profiles) == 0 {
		return nil, fmt.Errorf("no profiles found for uuid: %w", apperror.ErrNotFound)
	}

	contacts, err := s.contactSvc.GetForProfiles(ctx, profiles.GetIDs())
	if err != nil {
		return nil, err
	}

	contactsByEntity := contacts.MapByEntityID(entity.LinkedContactEntityProfileType)
	for i := 0; i < len(profiles); i++ {
		if profiles[i].LinkedUserAccess != "owner" {
			continue
		}

		if cont, ok := contactsByEntity[profiles[i].ID]; ok {
			profiles[i].Contacts = append(profiles[i].Contacts, cont)
		}
	}

	return profiles, nil
}
