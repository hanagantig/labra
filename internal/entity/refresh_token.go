package entity

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

type RefreshToken struct {
	token string
	hash  string
}

func NewRefreshToken() (RefreshToken, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return RefreshToken{}, err
	}

	token := base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(token))

	rt := RefreshToken{
		token: token,
		hash:  hex.EncodeToString(sum[:]),
	}

	return rt, nil
}

func NewRefreshTokenFromOpaque(opaque string) RefreshToken {
	sum := sha256.Sum256([]byte(opaque))

	rt := RefreshToken{
		token: opaque,
		hash:  hex.EncodeToString(sum[:]),
	}

	return rt
}

func NewHashedRefreshToken(hash string) RefreshToken {
	return RefreshToken{
		hash: hash,
	}
}

func (r RefreshToken) Hash() string {
	return r.hash
}

func (r RefreshToken) Plain() string {
	return r.token
}
