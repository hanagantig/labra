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

func TestAPI_VerifyUser(t *testing.T) {
	tests := map[string]struct {
		initialDBQueries []string
		requestBody      string

		expectedCode int
		expectedResp string

		checkData                      bool
		expectedProfiles               []expectedProfile
		expectedUsedCodes              []expectedCode
		expectedUsers                  []expectedUser
		expectedVerifiedLinkedContacts []expectedLinkedContact
		expectedUserProfiles           []expectedUserProfile
	}{
		"wrong request - return 400": {
			requestBody:  `unexpected`,
			expectedCode: 400,
			expectedResp: "{\"code\":400,\"error\":{\"message\":\"Unable to parse signup request: invalid character 'u' looking for beginning of value: bad request: client error\"}}",
		},
		"no contacts - return 404": {
			requestBody: `{
				"login": "me@test.loc",
				"password": "12345678910"
			}`,
			expectedCode: 404,
			expectedResp: "{\"code\":404,\"error\":{\"message\":\"Unable to verify user: get contact by value: item not found: client error\"}}",
		},
		"no linked contacts - return 404": {
			initialDBQueries: []string{
				"INSERT INTO users (id, uuid, f_name, l_name, password) VALUES (1, '00000000-0000000-0000000-0000001', 'test', 'test', '12345678910');",
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name)
					VALUES (5, '00000000-0000000-0000000-0000010', 1, 1, 'dfg', 'dfg');`,
				"INSERT INTO contacts (id, type, value) VALUES (555, 'email', 'me@test.loc');",
			},
			requestBody: `{
				"login": "me@test.loc",
				"password": "12345678910"
			}`,
			expectedCode: 404,
			expectedResp: "{\"code\":404,\"error\":{\"message\":\"Unable to verify user: no linked entities: item not found: client error\"}}",
		},
		"contact exists but no otp - return 404": {
			initialDBQueries: []string{
				"INSERT INTO users (id, uuid, f_name, l_name, password) VALUES (1, '00000000-0000000-0000000-0000001', 'test', 'test', '12345678910');",
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name)
					VALUES (5, '00000000-0000000-0000000-0000010', 1, 1, 'dfg', 'dfg');`,
				"INSERT INTO contacts (id, type, value) VALUES (555, 'email', 'me@test.loc');",
				"INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id) VALUES (1, 555, 'profile', '5');",
			},
			requestBody: `{
				"login": "me@test.loc",
				"password": "12345678910"
			}`,
			expectedCode: 404,
			expectedResp: "{\"code\":404,\"error\":{\"message\":\"Unable to verify user: get otp: item not found: client error\"}}",
		},
		"contact exists and otp with other type - return 404": {
			initialDBQueries: []string{
				"INSERT INTO users (id, uuid, f_name, l_name, password) VALUES (1, '00000000-0000000-0000000-0000001', 'test', 'test', '12345678910');",
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name)
					VALUES (5, '00000000-0000000-0000000-0000010', 1, 1, 'dfg', 'dfg');`,
				"INSERT INTO contacts (id, type, value) VALUES (555, 'email', 'me@test.loc');",
				"INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id) VALUES (1, 555, 'profile', '5');",
				"INSERT INTO codes (user_id, object_type, object_id, code) VALUES ('1', 'test', '1', '123');",
			},
			requestBody: `{
				"login": "me@test.loc",
				"password": "12345678910"
			}`,
			expectedCode: 404,
			expectedResp: "{\"code\":404,\"error\":{\"message\":\"Unable to verify user: get otp: item not found: client error\"}}",
		},
		"contact exists and otp for other user - return 404": {
			initialDBQueries: []string{
				"INSERT INTO users (id, uuid, f_name, l_name, password) VALUES (1, '00000000-0000000-0000000-0000001', 'test', 'test', '12345678910');",
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name)
					VALUES (5, '00000000-0000000-0000000-0000010', 1, 1, 'dfg', 'dfg');`,
				"INSERT INTO contacts (id, type, value) VALUES (555, 'email', 'me@test.loc');",
				"INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id) VALUES (1, 555, 'profile', '5');",
				"INSERT INTO codes (user_id, object_type, object_id, code) VALUES ('2', 'NEW_USER_CONTACT', '2', '123');",
			},
			requestBody: `{
				"login": "me@test.loc",
				"password": "12345678910"
			}`,
			expectedCode: 404,
			expectedResp: "{\"code\":404,\"error\":{\"message\":\"Unable to verify user: get otp: item not found: client error\"}}",
		},
		"otp exist with unexisting user - return 404 with error": {
			initialDBQueries: []string{
				"INSERT INTO users (id, uuid, f_name, l_name, password) VALUES (1, '00000000-0000000-0000000-0000001', 'test', 'test', '12345678910');",
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name)
					VALUES (5, '00000000-0000000-0000000-0000010', 1, 1, 'dfg', 'dfg');`,
				"INSERT INTO contacts (id, type, value) VALUES (555, 'email', 'me@test.loc');",
				"INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id) VALUES (1, 555, 'profile', '5');",
				"INSERT INTO codes (user_id, object_type, object_id, code, expired_at) VALUES ('2', 'NEW_USER_CONTACT', '555', '123', '2035-07-15 14:24:16');",
			},
			requestBody: `{
				"code":"123",
				"login": "me@test.loc",
				"password": "12345678910"
			}`,
			expectedCode: 404,
			expectedResp: "{\"code\":404,\"error\":{\"message\":\"Unable to verify user: invalid user id: item not found: client error\"}}",
		},
		//"contact with profile and otp exists - verify contact, create profile, use otp - return 200": {
		//	initialDBQueries: []string{
		//		"INSERT INTO users (id, uuid, f_name, l_name, password) VALUES (1, '00000000-0000000-0000000-0000001', 'test', 'test', '12345678910');",
		//		`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name)
		//			VALUES (5, '00000000-0000000-0000000-0000010', 1, 1, 'dfg', 'dfg');`,
		//		"INSERT INTO contacts (id, type, value) VALUES (555, 'email', 'me@test.loc');",
		//		"INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id) VALUES (1, 555, 'user', '1');",
		//		"INSERT INTO codes (user_id, object_type, object_id, code, expired_at) VALUES ('1', 'NEW_USER_CONTACT', '555', '123', '2035-07-15 14:24:16');",
		//	},
		//	requestBody: `{
		//		"code":"123",
		//		"login": "me@test.loc",
		//		"password": "12345678910"
		//	}`,
		//	expectedCode: 200,
		//	//expectedResp: "{\"code\":404,\"error\":{\"message\":\"Unable to verify user: get otp: item not found: client error\"}}",
		//
		//	checkData: true,
		//	expectedUsedCodes: []expectedCode{
		//		{
		//			UserID:     1,
		//			ObjectType: "NEW_USER_CONTACT",
		//			ObjectID:   "555",
		//			Code:       123,
		//		},
		//	},
		//	expectedProfiles: []expectedProfile{
		//		{
		//			ID:            5,
		//			Uuid:          "00000000-0000000-0000000-0000010",
		//			UserID:        1,
		//			CreatorUserID: 1,
		//		},
		//	},
		//	expectedVerifiedLinkedContacts: []expectedLinkedContact{
		//		{
		//			ID:         1,
		//			ContactID:  555,
		//			EntityType: "user",
		//			EntityID:   "1",
		//		},
		//	},
		//},
		"contact and otp exists, profile created by other account - verify contact, migrate profile, use otp - return 200": {
			initialDBQueries: []string{
				`INSERT INTO users (id, uuid, f_name, l_name, password) VALUES 
					(1, '00000000-0000000-0000000-0000001', 'test', 'test', '12345678910'),
					(2, '00000000-0000000-0000000-0000002', 'test', 'test', '12345678910');`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-0000000-0000000-0000010', 1, 1, 'dfg', 'dfg'),
				  	(5, '00000000-0000000-0000000-0000070', 0, 1, 'dfg', 'dfg')`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11, 1, 3, 'owner'),
				  	(22, 1, 5, 'owner');`,
				`INSERT INTO contacts (id, type, value) VALUES
					(100, 'email', 'first@test.loc'),
					(200, 'email', 'second@test.loc')`,
				`INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id) 
					VALUES (10, 100, 'user', '1'),
						   (20, 100, 'profile', '3'),
						   (30, 200, 'profile', '5'),
						   (40, 200, 'user', '2');`,
				`INSERT INTO codes (user_id, object_type, object_id, code, expired_at) 
					VALUES ('2', 'NEW_USER_CONTACT', '200', '123', '2035-07-15 14:24:16');`,
			},
			requestBody: `{
				"code":"123",
				"login": "second@test.loc",
				"password": "12345678910"
			}`,
			expectedCode: 200,

			checkData: true,
			expectedUsedCodes: []expectedCode{
				{
					UserID:     2,
					ObjectType: "NEW_USER_CONTACT",
					ObjectID:   "200",
					Code:       123,
				},
			},
			expectedProfiles: []expectedProfile{
				{
					ID:            3,
					UserID:        1,
					CreatorUserID: 1,
				},
				{
					ID:            5,
					UserID:        2,
					CreatorUserID: 1,
				},
			},
			expectedVerifiedLinkedContacts: []expectedLinkedContact{
				{
					ID:         30,
					ContactID:  200,
					EntityType: "profile",
					EntityID:   "5",
				},
				{
					ID:         40,
					ContactID:  200,
					EntityType: "user",
					EntityID:   "2",
				},
			},

			expectedUserProfiles: []expectedUserProfile{
				{
					UserID:      1,
					ProfileID:   3,
					AccessLevel: "owner",
				},
				{
					UserID:      1,
					ProfileID:   5,
					AccessLevel: "editor",
				},
				{
					UserID:      2,
					ProfileID:   5,
					AccessLevel: "owner",
				},
			},
		},
		"user with contact exists - create profile, use otp - return 200": {
			initialDBQueries: []string{
				`INSERT INTO users (id, uuid, f_name, l_name, password) VALUES 
					(1, '00000000-0000000-0000000-0000001', 'test', 'test', '12345678910'),
					(2, '00000000-0000000-0000000-0000002', 'test', 'test', '12345678910');`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-0000000-0000000-0000010', 1, 1, 'dfg', 'dfg');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11, 1, 3, 'owner');`,
				`INSERT INTO contacts (id, type, value) VALUES
					(100, 'email', 'first@test.loc'),
					(200, 'email', 'second@test.loc')`,
				`INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id) 
					VALUES (10, 100, 'user', '1'),
						   (20, 100, 'profile', '3'),
						   (30, 200, 'user', '2');`,
				`INSERT INTO codes (user_id, object_type, object_id, code, expired_at) 
					VALUES ('2', 'NEW_USER_CONTACT', '200', '123', '2035-07-15 14:24:16');`,
			},
			requestBody: `{
				"code":"123",
				"login": "second@test.loc",
				"password": "12345678910"
			}`,
			expectedCode: 200,

			checkData: true,
			expectedUsedCodes: []expectedCode{
				{
					UserID:     2,
					ObjectType: "NEW_USER_CONTACT",
					ObjectID:   "200",
					Code:       123,
				},
			},
			expectedProfiles: []expectedProfile{
				{
					ID:            3,
					UserID:        1,
					CreatorUserID: 1,
				},
				{
					ID:            4,
					UserID:        2,
					CreatorUserID: 2,
				},
			},
			expectedVerifiedLinkedContacts: []expectedLinkedContact{
				{
					ID:         30,
					ContactID:  200,
					EntityType: "user",
					EntityID:   "2",
				},
				{
					ID:         31,
					ContactID:  200,
					EntityType: "profile",
					EntityID:   "4",
				},
			},

			expectedUserProfiles: []expectedUserProfile{
				{
					UserID:      1,
					ProfileID:   3,
					AccessLevel: "owner",
				},
				{
					UserID:      2,
					ProfileID:   4,
					AccessLevel: "owner",
				},
			},
		},
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
				"http://localhost:8099/api/v1/user/verify",
				"application/json",
				bytes.NewBufferString(tc.requestBody),
			)
			require.NoError(t, err)

			body, _ := io.ReadAll(res.Body)
			defer res.Body.Close()

			require.Equal(t, tc.expectedCode, res.StatusCode, string(body))
			if tc.expectedCode != 200 {
				require.Equal(t, tc.expectedResp, string(body))
			}

			if !tc.checkData {
				return
			}

			var codes []expectedCode
			err = db.Select(&codes, `SELECT user_id, object_type, object_id, code FROM codes WHERE used_at IS NOT NULL;`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedUsedCodes, codes)

			var profiles []expectedProfile
			err = db.Select(&profiles, `SELECT id, user_id, creator_user_id FROM profiles`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedProfiles, profiles)

			var linkedContacts []expectedLinkedContact
			err = db.Select(&linkedContacts, `SELECT id, contact_id, entity_type, entity_id FROM linked_contacts WHERE verified_at IS NOT NULL;`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedVerifiedLinkedContacts, linkedContacts)

			var userProfiles []expectedUserProfile
			err = db.Select(&userProfiles, `SELECT user_id, profile_id, access_level FROM user_profiles;`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedUserProfiles, userProfiles)
		})
	}
}
