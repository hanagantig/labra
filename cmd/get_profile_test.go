package cmd

import (
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hanagantig/gracy"
	"github.com/stretchr/testify/require"
	"io"
	"labra/internal/app"
	"net/http"
	"testing"
	"time"
)

func TestAPI_GetProfile(t *testing.T) {
	a, err := app.NewApp("../config/app.conf.test.yaml")
	require.NoError(t, err)

	testUserUUID, _ := uuid.Parse("1bba9769-5f72-11f0-93b8-0242ac110002")

	tests := map[string]struct {
		initialDBQueries []string
		authToken        string

		expectedCode int
		expectedResp string
	}{
		"No token - return 401": {
			expectedCode: 401,
			expectedResp: "Invalid access token\n",
		},
		"Token with user id nil - return 401": {
			authToken:    newUserToken(uuid.Nil, a.GetConfig().Auth.AccessTokenSecret),
			expectedCode: 401,
			expectedResp: "invalid user uuid: Unauthorized: client error\n",
		},
		"No user - return 400": {
			authToken:    newUserToken(uuid.New(), a.GetConfig().Auth.AccessTokenSecret),
			expectedCode: 404,
			expectedResp: "{\"code\":404,\"error\":{\"message\":\"Unable to get user patients: no profiles found for uuid: item not found: client error\"}}",
		},
		"No profiles found - return 400": {
			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES (99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
			},
			authToken:    newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			expectedCode: 404,
			expectedResp: "{\"code\":404,\"error\":{\"message\":\"Unable to get user patients: no profiles found for uuid: item not found: client error\"}}",
		},
		"only one associated profile exists, no contacts for profile - return 200": {
			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES (99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner');`,
			},
			authToken:    newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			expectedCode: 200,
			expectedResp: `[
				{
				"access": "owner",
				"contacts": [],
				"date_of_birth": "0001-01-01 00:00:00 +0000 UTC",
				"id": "00000000-f000-0000-0000-000000000010",
				"patient_name": "Associated Profile"
			  }
			]`,
		},
		"only one associated profile exists with contacts for profile - return 200": {
			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES (99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner');`,
				`INSERT INTO contacts (id, type, value) VALUES
					(100, 'email', 'first@test.loc'),
					(200, 'email', 'second@test.loc');`,
				`INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id) 
					VALUES (10, 100, 'user', '99'),
						   (20, 100, 'profile', '3');`,
			},
			authToken:    newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			expectedCode: 200,
			expectedResp: `[
				{
				"access": "owner",
				"contacts": [
				  {
					"type": "email",
					"value": "first@test.loc"
				  }
				],
				"date_of_birth": "0001-01-01 00:00:00 +0000 UTC",
				"id": "00000000-f000-0000-0000-000000000010",
				"patient_name": "Associated Profile"
			  }
			]`,
		},
		"associated profile and added with contacts - return list of profiles - code 200": {
			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES (99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'owner');`,
				`INSERT INTO contacts (id, type, value) VALUES
					(100, 'email', 'first@test.loc'),
					(200, 'email', 'second@test.loc');`,
				`INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id) 
					VALUES (10, 100, 'user', '99'),
						   (20, 100, 'profile', '3'),
						   (30, 200, 'profile', '5');`,
			},
			authToken:    newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			expectedCode: 200,
			expectedResp: `[
			  {
				"access": "owner",
				"contacts": [
				  {
					"type": "email",
					"value": "first@test.loc"
				  }
				],
				"date_of_birth": "0001-01-01 00:00:00 +0000 UTC",
				"id": "00000000-f000-0000-0000-000000000010",
				"patient_name": "Associated Profile"
			  },
			  {
				"access": "owner",
				"contacts": [
				  {
					"type": "email",
					"value": "second@test.loc"
				  }
				],
				"date_of_birth": "0001-01-01 00:00:00 +0000 UTC",
				"id": "00000000-f000-0000-0000-000000000050",
				"patient_name": "Added Profile"
			  }
			]`,
		},
		"editor profiles with contacts - return list of profiles without contacts - code 200": {
			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES (99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'editor');`,
				`INSERT INTO contacts (id, type, value) VALUES
					(100, 'email', 'first@test.loc'),
					(200, 'email', 'second@test.loc');`,
				`INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id) 
					VALUES (10, 100, 'user', '99'),
						   (20, 100, 'profile', '3'),
						   (30, 200, 'profile', '5');`,
			},
			authToken:    newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			expectedCode: 200,
			expectedResp: `[
			  {
				"access": "owner",
				"contacts": [
				  {
					"type": "email",
					"value": "first@test.loc"
				  }
				],
				"date_of_birth": "0001-01-01 00:00:00 +0000 UTC",
				"id": "00000000-f000-0000-0000-000000000010",
				"patient_name": "Associated Profile"
			  },
			  {
				"access": "editor",
				"contacts": [],
				"date_of_birth": "0001-01-01 00:00:00 +0000 UTC",
				"id": "00000000-f000-0000-0000-000000000050",
				"patient_name": "Added Profile"
			  }
			]`,
		},
		"reader profiles with contacts - return list of profiles without contacts - code 200": {
			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES (99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'viewer');`,
				`INSERT INTO contacts (id, type, value) VALUES
					(100, 'email', 'first@test.loc'),
					(200, 'email', 'second@test.loc');`,
				`INSERT INTO linked_contacts (id, contact_id, entity_type, entity_id) 
					VALUES (10, 100, 'user', '99'),
						   (20, 100, 'profile', '3'),
						   (30, 200, 'profile', '5');`,
			},
			authToken:    newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			expectedCode: 200,
			expectedResp: `[
			  {
				"access": "owner",
				"contacts": [
				  {
					"type": "email",
					"value": "first@test.loc"
				  }
				],
				"date_of_birth": "0001-01-01 00:00:00 +0000 UTC",
				"id": "00000000-f000-0000-0000-000000000010",
				"patient_name": "Associated Profile"
			  },
			  {
				"access": "viewer",
				"contacts": [],
				"date_of_birth": "0001-01-01 00:00:00 +0000 UTC",
				"id": "00000000-f000-0000-0000-000000000050",
				"patient_name": "Added Profile"
			  }
			]`,
		},
	}

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
			req, err := http.NewRequest("GET", "http://localhost:8099/api/v1/account/profiles", nil)
			require.NoError(t, err)

			req.Header.Set("Authorization", "Bearer "+tc.authToken)
			req.Header.Set("Content-Type", "application/json")
			res, err := httpClient.Do(req)

			require.NoError(t, err)

			body, _ := io.ReadAll(res.Body)
			defer res.Body.Close()

			require.Equal(t, tc.expectedCode, res.StatusCode, string(body))
			assertJSONEqual(t, tc.expectedResp, string(body))
		})
	}
}
