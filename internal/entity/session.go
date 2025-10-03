package entity

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	AuthTokens AuthTokens
	SessionID  uuid.UUID
	UserUUID   uuid.UUID
	//Status    string
	//Binding   BindingType
	//JKT        *string // DPoP thumbprint
	DeviceID   string
	ExpiresAt  time.Time
	ReplacedAt time.Time
	ReplacedBy string
	RevokedAt  time.Time
}

func (s Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func (s Session) IsReplaced() bool {
	return !s.ReplacedAt.IsZero()
}

func (s Session) IsRevoked() bool {
	return !s.RevokedAt.IsZero()
}

func (s Session) IsActive() bool {
	return !s.IsExpired() && !s.IsReplaced() && !s.IsRevoked()
}

func (s Session) Replaced(newSession Session) Session {
	s.ReplacedAt = time.Now()
	s.ExpiresAt = time.Now()
	s.ReplacedBy = newSession.AuthTokens.RefreshToken.Hash()

	return s
}
