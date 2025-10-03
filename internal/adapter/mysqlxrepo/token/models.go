package token

import (
	"database/sql"
	"github.com/google/uuid"
	"labra/internal/entity"
	"time"
)

type RefreshToken struct {
	ID         int64        `db:"id"`
	UserUUID   string       `db:"user_uuid"`
	SessionID  string       `db:"session_id"`
	Token      string       `db:"token"`
	DeviceID   string       `db:"device_id"`
	ExpiresAt  sql.NullTime `db:"expires_at"`
	CreatedAt  time.Time    `db:"created_at"`
	ReplacedBy string       `db:"replaced_by"`
	ReplacedAt sql.NullTime `db:"replaced_at"`
	RevokedAt  sql.NullTime `db:"revoked_at"`
}

func NewRefreshTokenFromEntity(session entity.Session) RefreshToken {
	return RefreshToken{
		UserUUID:   session.UserUUID.String(),
		Token:      session.AuthTokens.RefreshToken.Hash(),
		SessionID:  session.SessionID.String(),
		DeviceID:   session.DeviceID,
		ReplacedBy: session.ReplacedBy,
		ReplacedAt: sql.NullTime{Time: session.ReplacedAt, Valid: !session.ReplacedAt.IsZero()},
		RevokedAt:  sql.NullTime{Time: session.RevokedAt, Valid: !session.RevokedAt.IsZero()},
		CreatedAt:  time.Now(),
		ExpiresAt:  sql.NullTime{Time: session.ExpiresAt, Valid: !session.ExpiresAt.IsZero()},
	}
}

func (r RefreshToken) ToEntity() entity.Session {
	userUUID, _ := uuid.Parse(r.UserUUID)
	sessionID, _ := uuid.Parse(r.SessionID)

	return entity.Session{
		AuthTokens: entity.AuthTokens{
			RefreshToken: entity.NewHashedRefreshToken(r.Token),
		},
		UserUUID:   userUUID,
		SessionID:  sessionID,
		DeviceID:   r.DeviceID,
		ReplacedBy: r.ReplacedBy,
		ReplacedAt: r.ReplacedAt.Time,
		ExpiresAt:  r.ExpiresAt.Time,
		RevokedAt:  r.RevokedAt.Time,
	}
}
