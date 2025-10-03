package user

import (
	"database/sql"
	"github.com/google/uuid"
	"labra/internal/entity"
)

type Users []User

type User struct {
	ID        int            `json:"id"`
	UUID      string         `db:"uuid"`
	FName     sql.NullString `db:"f_name"`
	LName     sql.NullString `db:"l_name"`
	Gender    sql.NullString `db:"gender"`
	BirthDate sql.NullTime   `db:"birth_date"`
	Password  string         `db:"password"`
}

func (u User) buildEntity() entity.User {
	userID, _ := uuid.Parse(u.UUID)

	return entity.User{
		ID:        u.ID,
		Uuid:      userID,
		FName:     u.FName.String,
		LName:     u.LName.String,
		BirthDate: u.BirthDate.Time,
		Password:  entity.UserPassword(u.Password),
	}
}

func NewUserFromEntity(u entity.User) User {
	return User{
		ID:        u.ID,
		UUID:      u.Uuid.String(),
		FName:     sql.NullString{String: u.FName, Valid: true},
		LName:     sql.NullString{String: u.LName, Valid: true},
		Password:  u.Password.String(),
		Gender:    sql.NullString{String: u.Gender, Valid: true},
		BirthDate: sql.NullTime{Time: u.BirthDate, Valid: true},
	}
}
