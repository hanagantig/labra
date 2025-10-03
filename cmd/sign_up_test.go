package cmd

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/hanagantig/gracy"
	"github.com/stretchr/testify/require"
	"io"
	"labra/internal/app"
	"net/http"
	"testing"
	"time"
)

func TestAPI_SignUp(t *testing.T) {
	tests := map[string]struct {
		initialDBQueries []string
		requestBody      string

		expectedStoredUsers    []expectedUser
		expectedStoredProfiles []expectedProfile
		expectedContacts       []expectedContact
		expectedLinkedContacts []expectedLinkedContact
		expectedCodes          []expectedCode

		expectedCode int
		expectedResp string
	}{
		"empty login - return 400 error": {
			requestBody: `{
				"login": "",
				"password": "12345678910"
			}`,
			expectedCode: 400,
			expectedResp: "{\"code\":400,\"error\":{\"message\":\"Invalid request: validation failure list:\\nlogin in body should be at least 5 chars long, bad request: client error\"}}",
		},
		"empty password - return 400 error": {
			requestBody: `{
				"login": "123456789",
				"password": ""
			}`,
			expectedCode: 400,
			expectedResp: "{\"code\":400,\"error\":{\"message\":\"Invalid request: validation failure list:\\npassword in body should be at least 10 chars long, bad request: client error\"}}",
		},
		"invalid email - return 400 error": {
			requestBody: `{
				"login": "123456789",
				"password": "12345567890"
			}`,
			expectedCode: 400,
			expectedResp: "{\"code\":400,\"error\":{\"message\":\"Invalid request: login should be email, bad request: client error\"}}",
		},
		// TODO: phone validation
		//"invalid phone - return 400 error": {},
		"verified contact exists and linked to an account - return 400 and ask to login": {
			initialDBQueries: []string{
				"INSERT INTO users (id, uuid, f_name, l_name, password) VALUES (1, '3d4faf99-60d1-11f0-93b8-0242ac110002', 'test', 'test', 'test');",
				"INSERT INTO contacts (id, type, value) VALUES (555, 'email', 'me@test.loc');",
				"INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id, verified_at) VALUES (3, 555, 'user', '3d4faf99-60d1-11f0-93b8-0242ac110002', '2025-07-14 19:44:22');",
			},
			requestBody: `{
				"login": "me@test.loc",
				"password": "12345567890"
			}`,

			expectedStoredUsers: []expectedUser{
				{
					ID:       1,
					Password: "test",
				},
			},

			expectedContacts: []expectedContact{
				{
					555,
					"email",
					"me@test.loc",
				},
			},

			expectedLinkedContacts: []expectedLinkedContact{
				{
					ID:         3,
					ContactID:  555,
					EntityType: "user",
					EntityID:   "3d4faf99-60d1-11f0-93b8-0242ac110002",
					VerifiedAt: "2025-07-14 19:44:22",
				},
			},
			expectedCode: 400,
			expectedResp: "{\"code\":400,\"error\":{\"message\":\"Unable to sign up: entity already exists: client error\"}}",
		},
		"unverified contact exists and linked to an account - return 400 and ask to login": {
			initialDBQueries: []string{
				"INSERT INTO users (id, uuid, f_name, l_name, password) VALUES (1, '3d4faf99-60d1-11f0-93b8-0242ac110002', 'test', 'test', 'test');",
				"INSERT INTO contacts (id, type, value) VALUES (555, 'email', 'me@test.loc');",
				"INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id) VALUES (3, 555, 'user', '3d4faf99-60d1-11f0-93b8-0242ac110002');",
			},
			requestBody: `{
				"login": "me@test.loc",
				"password": "12345567890"
			}`,

			expectedStoredUsers: []expectedUser{
				{
					ID:       1,
					Password: "test",
				},
			},

			expectedContacts: []expectedContact{
				{
					555,
					"email",
					"me@test.loc",
				},
			},

			expectedLinkedContacts: []expectedLinkedContact{
				{
					ID:         3,
					ContactID:  555,
					EntityType: "user",
					EntityID:   "3d4faf99-60d1-11f0-93b8-0242ac110002",
				},
			},

			expectedCode: 400,
			expectedResp: "{\"code\":400,\"error\":{\"message\":\"Unable to sign up: entity already exists: client error\"}}",
		},
		"verified contact exists and linked to a profile  - create account, otp generated verification, return 200": {
			initialDBQueries: []string{
				"INSERT INTO users (id, uuid, f_name, l_name, password) VALUES (1, '00000000-0000000-0000000-0000001', 'test', 'test', 'test');",
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name)
					VALUES (5, '00000000-0000000-0000000-0000010', 1, 1, 'dfg', 'dfg');`,
				"INSERT INTO contacts (id, type, value) VALUES (555, 'email', 'me@test.loc');",
				"INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id, verified_at) VALUES (1, 555, 'profile', '5', '2025-07-14 20:15:31');",
			},
			requestBody: `{
				"login": "me@test.loc",
				"password": "12345567890"
			}`,

			expectedStoredUsers: []expectedUser{
				{
					ID:       1,
					Password: "test",
				},
				{
					ID:       2,
					Password: "12345567890",
				},
			},
			expectedStoredProfiles: []expectedProfile{
				{
					ID:            5,
					Uuid:          "00000000-0000000-0000000-0000010",
					UserID:        1,
					CreatorUserID: 1,
				},
			},

			expectedContacts: []expectedContact{
				{
					ID:    555,
					Type:  "email",
					Value: "me@test.loc",
				},
			},
			expectedLinkedContacts: []expectedLinkedContact{
				{
					ID:         1,
					ContactID:  555,
					EntityType: "profile",
					EntityID:   "5",
					VerifiedAt: "2025-07-14 20:15:31",
				},
				{
					ID:         2,
					ContactID:  555,
					EntityType: "user",
					EntityID:   "2",
				},
			},

			expectedCodes: []expectedCode{
				{
					UserID:     2,
					ObjectType: "NEW_USER_CONTACT",
					ObjectID:   "555",
				},
			},

			expectedCode: 200,
			expectedResp: "{\"success\":\"OK\"}",
		},
		"unverified contact exists and linked to a profile - create account, ask for verification, return 200": {
			initialDBQueries: []string{
				"INSERT INTO users (id, uuid, f_name, l_name, password) VALUES (1, '00000000-0000000-0000000-0000001', 'test', 'test', 'test');",
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name)
					VALUES (5, '00000000-0000000-0000000-0000010', 1, 1, 'dfg', 'dfg');`,
				"INSERT INTO contacts (id, type, value) VALUES (555, 'email', 'me@test.loc');",
				"INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id) VALUES (1, 555, 'profile', '5');",
			},
			requestBody: `{
				"login": "me@test.loc",
				"password": "12345567890"
			}`,

			expectedStoredUsers: []expectedUser{
				{
					ID:       1,
					Password: "test",
				},
				{
					ID:       2,
					Password: "12345567890",
				},
			},
			expectedStoredProfiles: []expectedProfile{
				{
					ID:            5,
					Uuid:          "00000000-0000000-0000000-0000010",
					UserID:        1,
					CreatorUserID: 1,
				},
			},

			expectedContacts: []expectedContact{
				{
					ID:    555,
					Type:  "email",
					Value: "me@test.loc",
				},
			},
			expectedLinkedContacts: []expectedLinkedContact{
				{
					ID:         1,
					ContactID:  555,
					EntityType: "profile",
					EntityID:   "5",
				},
				{
					ID:         2,
					ContactID:  555,
					EntityType: "user",
					EntityID:   "2",
				},
			},

			expectedCodes: []expectedCode{
				{
					UserID:     2,
					ObjectType: "NEW_USER_CONTACT",
					ObjectID:   "555",
				},
			},

			expectedCode: 200,
			expectedResp: "{\"success\":\"OK\"}",
		},
		"no accounts and no contact - user successfully signed up by email, return 200": {
			requestBody: `{
				"login": "me@test.loc",
				"password": "12345567890"
			}`,

			expectedStoredUsers: []expectedUser{
				{
					ID:       1,
					Password: "12345567890",
				},
			},

			expectedContacts: []expectedContact{
				{
					ID:    1,
					Type:  "email",
					Value: "me@test.loc",
				},
			},
			expectedLinkedContacts: []expectedLinkedContact{
				{
					ID:         1,
					ContactID:  1,
					EntityType: "user",
					EntityID:   "1",
				},
			},

			expectedCodes: []expectedCode{
				{
					UserID:     1,
					ObjectType: "NEW_USER_CONTACT",
					ObjectID:   "1",
				},
			},

			expectedCode: 200,
			expectedResp: "{\"success\":\"OK\"}",
		},
		//"no accounts and no contact - user successfully signed up by phone, return 200":       {},
		//"Sign-up with same contact multiple times in short window - rate limit error":         {},
		//"Mass sign-ups with malformed or fake emails/phones (spam attack) - rate limit error": {},
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

			// Clean DB
			truncateAllTables(db, "mego_api_test")
			for _, query := range tc.initialDBQueries {
				_ = db.MustExec(query)
			}
			defer truncateAllTables(db, "mego_api_test")

			httpClient := &http.Client{}
			res, err := httpClient.Post(
				"http://localhost:8099/api/v1/signup",
				"application/json",
				bytes.NewBufferString(tc.requestBody),
			)
			require.NoError(t, err)

			body, _ := io.ReadAll(res.Body)
			defer res.Body.Close()

			require.Equal(t, tc.expectedCode, res.StatusCode, string(body))
			require.Equal(t, tc.expectedResp, string(body))

			var users []expectedUser
			userQuery := `SELECT id, password FROM users`

			err = db.Select(&users, userQuery)
			require.NoError(t, err)
			require.Equal(t, tc.expectedStoredUsers, users)

			var profiles []expectedProfile
			err = db.Select(&profiles, `SELECT id, uuid, user_id, creator_user_id FROM profiles`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedStoredProfiles, profiles)

			var contacts []expectedContact
			err = db.Select(&contacts, `SELECT id, type, value FROM contacts`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedContacts, contacts)

			var linkedContacts []expectedLinkedContact
			err = db.Select(&linkedContacts, `SELECT id, contact_id, entity_type, entity_id, COALESCE(verified_at, '') AS verified_at FROM linked_contacts`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedLinkedContacts, linkedContacts)

			var codes []expectedCode
			err = db.Select(&codes, `SELECT user_id, object_type, object_id FROM codes`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedCodes, codes)
		})
	}
}
