package profile

import (
	"database/sql"
	"github.com/google/uuid"
	"labra/internal/entity"
)

type Profiles []Profile
type Profile struct {
	ID            int            `db:"id"`
	UUID          string         `db:"uuid"`
	UserID        sql.NullInt64  `db:"user_id"`
	LinkedUserID  sql.NullInt64  `db:"linked_user_id"`
	AccessLevel   string         `db:"access_level"`
	CreatorUserID int            `db:"creator_user_id"`
	FName         sql.NullString `db:"f_name"`
	LName         sql.NullString `db:"l_name"`
	Gender        sql.NullString `db:"gender"`
	BirthDate     sql.NullTime   `db:"birth_date"`
}

func NewFromEntity(e entity.Profile) Profile {
	bDate := sql.NullTime{}
	if !e.DateOfBirth.IsZero() {
		bDate = sql.NullTime{Time: e.DateOfBirth, Valid: true}
	}

	return Profile{
		ID:            e.ID,
		UUID:          e.Uuid.String(),
		UserID:        sql.NullInt64{Valid: e.UserID > 0, Int64: int64(e.UserID)},
		CreatorUserID: e.CreatorUserID,
		FName:         stringToNullString(e.FName),
		LName:         stringToNullString(e.LName),
		Gender:        stringToNullString(e.Gender),
		BirthDate:     bDate,
	}
}

func stringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}

	return sql.NullString{Valid: true, String: s}
}

func (p Profile) buildEntity() entity.Profile {
	pUUID, err := uuid.Parse(p.UUID)
	_ = err

	return entity.Profile{
		ID:               p.ID,
		Uuid:             pUUID,
		UserID:           int(p.UserID.Int64),
		CreatorUserID:    p.CreatorUserID,
		FName:            p.FName.String,
		LName:            p.LName.String,
		Gender:           p.Gender.String,
		DateOfBirth:      p.BirthDate.Time,
		LinkedUserAccess: p.AccessLevel,
		LinkedUserID:     int(p.LinkedUserID.Int64),
	}
}

func (p Profiles) buildEntity() entity.Profiles {
	res := make(entity.Profiles, 0, len(p))
	for _, patient := range p {
		res = append(res, patient.buildEntity())
	}

	return res
}
