package cmd

import (
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hanagantig/gracy"
	"github.com/stretchr/testify/require"
	"io"
	"labra/internal/app"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestAPI_PostProfile(t *testing.T) {
	a, err := app.NewApp("../config/app.conf.test.yaml")
	require.NoError(t, err)

	testUserUUID, _ := uuid.Parse("1bba9769-5f72-11f0-93b8-0242ac110002")

	tests := map[string]struct {
		initialDBQueries []string
		requestBody      string
		authToken        string

		expectedCode int
		expectedResp string

		expectedProfiles       []expectedProfile
		expectedUserProfiles   []expectedUserProfile
		expectedContacts       []expectedContact
		expectedLinkedContacts []expectedLinkedContact
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
		"Invalid request - return 400": {
			authToken:    newUserToken(uuid.New(), a.GetConfig().Auth.AccessTokenSecret),
			requestBody:  `invalid`,
			expectedCode: 400,
			expectedResp: "{\"code\":400,\"error\":{\"message\":\"Unable to parse add profile request: invalid character 'i' looking for beginning of value, bad request: client error\"}}",
		},
		"Wrong request parameters - return 404": {
			authToken: newUserToken(uuid.New(), a.GetConfig().Auth.AccessTokenSecret),
			requestBody: `{
    			"genderrp":"M",
    			"patient_namee":"test"
			}`,
			expectedCode: 400,
			expectedResp: "{\"code\":400,\"error\":{\"message\":\"Invalid request: validation failure list:\\nfirst_name in body is required\\ngender in body is required, bad request: client error\"}}",
		},
		"No user with uuid - return 404": {
			authToken: newUserToken(uuid.New(), a.GetConfig().Auth.AccessTokenSecret),
			requestBody: `{
    			"gender":"M",
    			"first_name":"test"
			}`,
			expectedCode: 404,
			expectedResp: "{\"code\":404,\"error\":{\"message\":\"Unable to create a profile: can't get user by uuid: item not found: client error\"}}",
		},
		"User doesn't have any profile - return 404": {
			authToken: newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			initialDBQueries: []string{
				"INSERT INTO users(id, uuid, l_name, password) VALUES (99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')",
			},
			requestBody: `{
    			"gender":"M",
    			"first_name":"test"
			}`,
			expectedCode: 404,
			expectedResp: "{\"code\":404,\"error\":{\"message\":\"Unable to create a profile: can't get user profiles: no profiles found for uuid: item not found: client error\"}}",
		},
		"User has associated profile - create none associated one - return 200": {
			authToken: newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			initialDBQueries: []string{
				"INSERT INTO users(id, uuid, l_name, password) VALUES (99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')",
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, gender) VALUES 
				  	(3, '00000000-0000000-0000000-0000010', 99, 99, 'dfg', 'F');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES
					(55, 99, 3, "owner");`,
			},
			requestBody: `{
    			"gender":"M",
    			"first_name":"test"
			}`,
			expectedCode: 200,

			expectedProfiles: []expectedProfile{
				{
					ID:            3,
					UserID:        99,
					CreatorUserID: 99,
					Gender:        "F",
					FName:         "dfg",
				},
				{
					ID:            4,
					UserID:        0,
					CreatorUserID: 99,
					Gender:        "M",
					FName:         "test",
				},
			},
			expectedUserProfiles: []expectedUserProfile{
				{
					UserID:      99,
					ProfileID:   3,
					AccessLevel: "owner",
				},
				{
					UserID:      99,
					ProfileID:   4,
					AccessLevel: "owner",
				},
			},
		},
		"User has multiple profiles - create none associated one - return 200": {
			authToken: newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			initialDBQueries: []string{
				"INSERT INTO users(id, uuid, l_name, password) VALUES (99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')",
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, gender) VALUES 
				  	(3, '00000000-0000000-0000000-0000010', 99, 99, 'dfg', 'F'),
				  	(4, '00000000-0000000-0000000-0000050', null, 99, 'profile2', 'M');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES
					(55, 99, 3, "owner"),
					(56, 99, 4, "editor");`,
			},
			requestBody: `{
    			"gender":"F",
    			"first_name":"test"
			}`,
			expectedCode: 200,

			expectedProfiles: []expectedProfile{
				{
					ID:            3,
					UserID:        99,
					CreatorUserID: 99,
					Gender:        "F",
					FName:         "dfg",
				},
				{
					ID:            4,
					UserID:        0,
					CreatorUserID: 99,
					Gender:        "M",
					FName:         "profile2",
				},
				{
					ID:            5,
					UserID:        0,
					CreatorUserID: 99,
					Gender:        "F",
					FName:         "test",
				},
			},
			expectedUserProfiles: []expectedUserProfile{
				{
					UserID:      99,
					ProfileID:   3,
					AccessLevel: "owner",
				},
				{
					UserID:      99,
					ProfileID:   4,
					AccessLevel: "editor",
				},
				{
					UserID:      99,
					ProfileID:   5,
					AccessLevel: "owner",
				},
			},
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
			req, err := http.NewRequest("POST", "http://localhost:8099/api/v1/account/profiles", strings.NewReader(tc.requestBody))
			require.NoError(t, err)

			req.Header.Set("Authorization", "Bearer "+tc.authToken)
			req.Header.Set("Content-Type", "application/json")
			res, err := httpClient.Do(req)

			require.NoError(t, err)

			body, _ := io.ReadAll(res.Body)
			defer res.Body.Close()

			require.Equal(t, tc.expectedCode, res.StatusCode, string(body))
			if tc.expectedCode != 200 {
				require.Equal(t, tc.expectedResp, string(body))
			}

			var profiles []expectedProfile
			err = db.Select(&profiles, `SELECT id, coalesce(user_id,0) as user_id, creator_user_id, f_name, gender FROM profiles`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedProfiles, profiles)

			var userProfiles []expectedUserProfile
			err = db.Select(&userProfiles, `SELECT user_id, profile_id, access_level FROM user_profiles;`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedUserProfiles, userProfiles)

			var linkedContacts []expectedLinkedContact
			err = db.Select(&linkedContacts, `SELECT id, contact_id, entity_type, entity_id FROM linked_contacts WHERE verified_at IS NOT NULL;`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedLinkedContacts, linkedContacts)
		})
	}
}
