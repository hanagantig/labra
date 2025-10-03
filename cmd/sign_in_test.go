package cmd

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/hanagantig/gracy"
	"github.com/stretchr/testify/require"
	"io"
	"labra/internal/app"
	"labra/internal/entity"
	"labra/internal/handler/http/api/v1/models"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestAPI_SignIn(t *testing.T) {
	tests := map[string]struct {
		initialDBQueries []string
		requestBody      string

		expectedCode     int
		expectedUserUUID string
		expectedUserID   int
		expectedErr      string
	}{
		"user does not exists - return 404 error": {
			requestBody: `{
				"login": "login",
    			"password": "pass"	
			}`,

			expectedCode: http.StatusNotFound,
			expectedErr:  `{"code":404,"error":{"message":"item not found: client error"}}`,
		},
		"user exists but password is wrong - return 404 error": {
			initialDBQueries: []string{
				"INSERT INTO users(id, uuid, l_name, password) VALUES (5, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')",
				"INSERT INTO contacts (id, type, value) VALUES (1, 'email', 'test')",
				"INSERT INTO linked_contacts (contact_id, entity_type, entity_id) VALUES (1, 'user', 5)",
			},

			requestBody: `{
				"login": "login",
    			"password": "pass"	
			}`,

			expectedCode: http.StatusNotFound,
			expectedErr:  `{"code":404,"error":{"message":"item not found: client error"}}`,
		},
		"user exists and password is correct - return response": {
			initialDBQueries: []string{
				"INSERT INTO users(id, uuid, l_name, password) VALUES (5, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', '$2a$10$uK2/QDBNMsAx8glKYG6wo.8yCjhx5hDwWHuFhH4.P.lwwqZoWsPsC')",
				"INSERT INTO mego_api_test.contacts (id, type, value) VALUES (1, 'email', 'test')",
				"INSERT INTO linked_contacts (contact_id, entity_type, entity_id) VALUES (1, 'user', 5)",
			},

			requestBody: `{
				"login": "test",
    			"password": "pass"
			}`,

			expectedCode:     http.StatusOK,
			expectedUserUUID: "1bba9769-5f72-11f0-93b8-0242ac110002",
			expectedUserID:   5,
		},
		//"user sign in with unverified contact - return verification request":                                                 {},
		//"too many attempts to login - return restriction error":                                                              {},
		//"Missing required login → return 400 Bad Request":                                                                    {},
		//"Missing required password → return 400 Bad Request":                                                                 {},
		//"Malformed email or identifier format → return 400 Bad Request":                                                      {},
		//"Empty string or whitespace-only password → return 400 Bad Request":                                                  {},
		//"User exists, account is disabled → return 403 Forbidden":                                                            {},
		//"User exists but marked as deleted → return 404 Not Found":                                                           {},
		//"Correct password, but device or IP not trusted → return 403 Forbidden + trigger 2FA or email challenge":             {},
		//"Correct credentials, but MFA is enabled and not yet passed → return MFA challenge required":                         {},
		//"Multiple failed logins from same IP → apply progressive delay or lockout mechanism":                                 {},
		//"Multiple failed logins across IPs for same user → temporary account lock or challenge":                              {},
		//"Timing attack protection (uniform response time) → measure response time to ensure no leak":                         {},
		//"Successful login issues token with proper expiry and scopes → verify returned token format and claims":              {},
		//"User signs in with uppercase/lowercase variation of email → verify case sensitivity policy":                         {},
		//"User signs in during server clock drift or NTP issue → ensure timestamps are handled properly":                      {},
		//"Username collision with reserved or system keywords (e.g., \"admin\", \"root\") → ensure blocked or treated safely": {},
		//"Injected SQL or XSS in email/password fields → ensure safely rejected / escaped":                                    {},
		//"Long inputs (e.g., 1000+ chars in password/email) → return 400 Bad Request, test for DoS vector":                    {},
		//"Non-UTF8 characters in input → return 400 Bad Request":                                                              {},
		//"": {},
	}

	a, err := app.NewApp("../config/app.conf.test.yaml")
	require.NoError(t, err)

	go a.StartHTTPServer()
	defer gracy.GracefulShutdown()

	db := a.GetDB()

	time.Sleep(1 * time.Second)
	for testName, tc := range tests {
		tc := tc
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			truncateAllTables(db, "mego_api_test")
			for _, query := range tc.initialDBQueries {
				_ = db.MustExec(query)
			}
			defer truncateAllTables(db, "mego_api_test")

			httpClient := &http.Client{}
			res, err := httpClient.Post(
				"http://localhost:8099/api/v1/sign_in",
				"application/json",
				bytes.NewBufferString(tc.requestBody),
			)
			require.NoError(t, err)

			body, _ := io.ReadAll(res.Body)
			defer res.Body.Close()

			require.Equal(t, tc.expectedCode, res.StatusCode)
			if tc.expectedErr != "" {
				require.Equal(t, tc.expectedErr, string(body))
				return
			}

			var resp models.APIAuthToken
			decoder := json.NewDecoder(strings.NewReader(string(body)))
			err = decoder.Decode(&resp)
			require.NoError(t, err, "Failed to decode JSON")

			// Validate access token
			accessToken := entity.JWT(resp.AccessToken)
			claims, err := accessToken.ValidateAndGetClientClaims(a.GetConfig().Auth.AccessTokenSecret)
			require.NoError(t, err, "Invalid JWT token")

			// Check required claims
			require.NotEmpty(t, claims.Subject, "sub claim must be present")
			require.Equal(t, tc.expectedUserUUID, claims.Subject)
			require.Equal(t, "API", claims.Issuer, "unexpected issuer")

			// Check iat and nbf
			now := time.Now().Unix()
			require.LessOrEqual(t, claims.IssuedAt.Unix(), now, "iat must be <= now")
			require.LessOrEqual(t, claims.NotBefore.Unix(), now, "nbf must be <= now")
			require.NotNil(t, claims.ExpiresAt)
			require.Greater(t, claims.ExpiresAt.Unix(), now, "exp must be > now")
			require.Equal(t, 1800, int(claims.ExpiresAt.Unix()-claims.IssuedAt.Unix()))

			// Validate refresh token format (basic check)
			require.NotEmpty(t, resp.RefreshToken, "refresh token must be present")

			var refreshRes struct {
				RefreshToken string       `db:"token"`
				ExpiresAt    time.Time    `db:"expires_at"`
				CreatedAt    time.Time    `db:"created_at"`
				RevokedAt    sql.NullTime `db:"revoked_at"`
				UserUUID     string       `db:"user_uuid"`
			}

			refreshToken := entity.NewRefreshTokenFromOpaque(resp.RefreshToken)
			err = db.Get(&refreshRes, "SELECT token, expires_at, user_uuid, created_at FROM sessions WHERE token = ? LIMIT 1", refreshToken.Hash())
			require.NoError(t, err)

			diff := (refreshRes.ExpiresAt.Unix() - refreshRes.CreatedAt.Unix()) - int64(a.GetConfig().Auth.RefreshTokenTTL.Seconds())
			require.LessOrEqual(t, int(diff), 5)
			require.Equal(t, time.Time{}, refreshRes.RevokedAt.Time)
			require.Equal(t, tc.expectedUserUUID, refreshRes.UserUUID)

		})
	}
}
