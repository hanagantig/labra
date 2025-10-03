package middleware

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"labra/internal/entity"
	"net/http"
	"strings"
)

type ClientID string

type ClientPhone string

const UserIDKey = ClientID("user_id")

const UserPatientsKey = ClientPhone("user_patients")

func AuthorizationMiddleware(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Обработка CORS-заголовков
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			//w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")

			authHeader := r.Header.Get("Authorization")
			t := strings.Split(authHeader, " ")
			if len(t) != 2 {
				http.Error(w, "Invalid access token", http.StatusUnauthorized)
				return
			}

			authToken := entity.JWT(t[1])
			claims, err := authToken.ValidateAndGetClientClaims(secret)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, claims.Subject)
			//ctx = context.WithValue(ctx, UserPatientsKey, claims.UserPatients)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserUUIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return uuid.Nil, errors.New("user UUID not found in context")
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}
