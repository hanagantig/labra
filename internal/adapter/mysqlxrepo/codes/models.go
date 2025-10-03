package codes

import (
	"labra/internal/entity"
	"time"
)

type Code struct {
	ID         uint      `db:"id"`
	UserID     int       `db:"user_id"`
	ObjectType string    `db:"object_type"`
	ObjectID   string    `db:"object_id"`
	Code       string    `db:"code"`
	ExpiredAt  time.Time `db:"expired_at"`
	CreatedAt  time.Time `db:"created_at"`
}

func NewCodeFromEntity(e entity.OTPCode) Code {
	return Code{
		ID:         e.ID,
		UserID:     e.UserID,
		ObjectType: e.ObjectType.String(),
		ObjectID:   e.ObjectID,
		Code:       e.Code,
		ExpiredAt:  e.ExpiredAt,
		CreatedAt:  e.CreatedAt,
	}
}

func (c Code) buildEntity() entity.OTPCode {
	return entity.OTPCode{
		ID:         c.ID,
		UserID:     c.UserID,
		ObjectType: entity.OTPObjectType(c.ObjectType),
		ObjectID:   c.ObjectID,
		Code:       c.Code,
		ExpiredAt:  c.ExpiredAt,
		CreatedAt:  c.CreatedAt,
	}
}
