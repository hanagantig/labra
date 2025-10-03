package user

import (
	"context"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (s *Service) GetAuthUser(ctx context.Context, login entity.EmailOrPhone, pass entity.UserPassword) (entity.UserWithProfiles, error) {
	user, err := s.userRepo.GetByCredentials(ctx, login, pass)
	if err != nil {
		return entity.UserWithProfiles{}, apperror.ToAppError(err)
	}
	//TODO: check if user verified
	return entity.UserWithProfiles{
		User: user,
		//Profiles: patients,
	}, nil
}
