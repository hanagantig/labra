package usecase

import (
	"context"
	"labra/internal/entity"
)

func (u *UseCase) UserSignUp(ctx context.Context, login entity.EmailOrPhone, pass entity.UserPassword) (entity.User, error) {
	hashedPass, err := pass.Hashed()
	if err != nil {
		return entity.User{}, err
	}

	user := entity.User{
		Password: hashedPass,
		//TODO: add more details (name, birth_date ...)
	}

	// TODO: handle device !!!

	return u.userSvc.CreateUser(ctx, login, user)
}
