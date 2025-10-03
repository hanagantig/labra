package entity

import (
	"fmt"
	"github.com/google/uuid"
	"labra/internal/apperror"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type AuthTokens struct {
	AccessToken  JWT
	RefreshToken RefreshToken
}

type JWT string

type ClientClaims struct {
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

func NewUserJWT(s Session, secret string, ttl time.Duration) (JWT, error) {
	claims := ClientClaims{
		s.SessionID.String(),
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "API",
			Subject:   s.UserUUID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return JWT(t), nil
}

func (j JWT) String() string {
	return string(j)
}

func (j JWT) ValidateAndGetClientClaims(secret string) (ClientClaims, error) {
	token, err := jwt.ParseWithClaims(j.String(), &ClientClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return ClientClaims{}, err
	}

	claims, ok := token.Claims.(*ClientClaims)
	if !ok || !token.Valid {
		return ClientClaims{}, fmt.Errorf("invalid token")
	}

	userUuid, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ClientClaims{}, err
	}

	if userUuid == uuid.Nil {
		return ClientClaims{}, fmt.Errorf("invalid user uuid: %w", apperror.ErrUnauthorized)
	}

	return *claims, nil
}
