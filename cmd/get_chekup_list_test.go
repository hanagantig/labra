package cmd

import (
	"fmt"
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

func TestAPI_GetCheckupList(t *testing.T) {
	a, err := app.NewApp("../config/app.conf.test.yaml")
	require.NoError(t, err)

	testUserUUID, _ := uuid.Parse("1bba9769-5f72-11f0-93b8-0242ac110002")

	tests := map[string]struct {
		initialDBQueries []string
		authToken        string
		profileUUID      string
		queryString      string

		expectedCode int
		expectedResp string
	}{
		"No profile id in url - return 404": {
			expectedCode: 404,
			expectedResp: "404 page not found\n",
		},
		"No token - return 401": {
			profileUUID:  "test",
			expectedCode: 401,
			expectedResp: "Invalid access token\n",
		},
		"Token with user id nil - return 401": {
			profileUUID:  "test",
			authToken:    newUserToken(uuid.Nil, a.GetConfig().Auth.AccessTokenSecret),
			expectedCode: 401,
			expectedResp: "invalid user uuid: Unauthorized: client error\n",
		},
		"Invalid profile id - return 400": {
			profileUUID:  "test",
			authToken:    newUserToken(uuid.New(), a.GetConfig().Auth.AccessTokenSecret),
			expectedCode: 400,
			expectedResp: ` {
				"code":400, 
				"error": {
					"message":"profile id is required: invalid UUID length: 4, bad request: client error"
				}
			}`,
		},
		"No user - return 404": {
			profileUUID:  uuid.New().String(),
			authToken:    newUserToken(uuid.New(), a.GetConfig().Auth.AccessTokenSecret),
			expectedCode: 404,
			expectedResp: `{
				"code":404, 
				"error": {
					"message":"unable to get list of checkups: can't get user by uuid: item not found: client error"
				}
			}`,
		},
		"No user profile - return 404": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),

			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES 
				  	(99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
			},
			expectedCode: 404,
			expectedResp: `{
				"code":404, 
				"error": {
					"message":"unable to get list of checkups: can't get user profiles: no profiles found for uuid: item not found: client error"
				}
			}`,
		},
		"No profiles - return 404": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),

			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES 
				  	(99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'owner');`,
			},
			expectedCode: 404,
			expectedResp: `{
				"code":404, 
				"error": {
					"message":"unable to get list of checkups: can't get user profiles: no profiles found for uuid: item not found: client error"
				}
			}`,
		},
		"Profile doesn't have checkups - return empty response": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),

			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES 
				  	(99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'owner');`,
			},
			expectedCode: 200,
			expectedResp: `[]`,
		},
		"Checkup doesn't have results - return empty response": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),

			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES 
				  	(99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'owner');`,
				`INSERT INTO checkups (id, profile_id, lab_id, status, date) 
					VALUES (1, 3, 0, 'TEST', '2025-07-18');`,
			},
			expectedCode: 200,
			expectedResp: `[]`,
		},
		"Checkup with results - return response with results": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),

			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES 
				  	(99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'owner');`,
				`INSERT INTO checkups (id, profile_id, lab_id, status, date) 
					VALUES (10, 3, 0, 'TEST', '2025-07-18');`,
				`INSERT INTO checkup_results (checkup_id, marker_id, unit_id, value) VALUES 
			  		(10, 55, 33, 100),
			  		(10, 77, 11, 0.53);`,
			},
			expectedCode: 200,
			expectedResp: `[ 
				{
				  "checkup" : {
					"date" : "2025-07-18 00:00:00 +0000 UTC",
					"id" : 10,
					"lab" : {
					  "id" : 0,
					  "name" : ""
					},
					"profile" : {
					  "id" : "00000000-f000-0000-0000-000000000010",
					  "first_name" : "Associated"
					}
				  },
				  "results" : [ 
					{
						"id" : 55,
						"name" : "",
						"unit" : {
						  "id" : 33,
						  "name" : ""
						},
						"value" : 100
				  	},
					{
						"id" : 77,
						"name" : "",
						"unit" : {
						  "id" : 11,
						  "name" : ""
						},
						"value" : 0.53
				  	}
				  ]
				} 
			]`,
		},
		"Checkup with results and markers - return response with results": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),

			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES 
				  	(99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'owner');`,
				`INSERT INTO checkups (id, profile_id, lab_id, status, date) 
					VALUES (10, 3, 0, 'TEST', '2025-07-18');`,
				`INSERT INTO checkup_results (checkup_id, marker_id, unit_id, value) VALUES 
			  		(10, 55, 33, 100),
			  		(10, 77, 11, 0.53);`,
				`INSERT INTO markers (id, name, ref_range_min, ref_range_max) VALUES 
			 		(55, 'M1', 10, 50),
			 		(77, 'M2', 10, 50);`,
			},
			expectedCode: 200,
			expectedResp: `[ 
				{
				  "checkup" : {
					"date" : "2025-07-18 00:00:00 +0000 UTC",
					"id" : 10,
					"lab" : {
					  "id" : 0,
					  "name" : ""
					},
					"profile" : {
					  "id" : "00000000-f000-0000-0000-000000000010",
					  "first_name" : "Associated"
					}
				  },
				  "results" : [ 
					{
						"id" : 55,
						"name" : "M1",
						"unit" : {
						  "id" : 33,
						  "name" : ""
						},
						"value" : 100
				  	},
					{
						"id" : 77,
						"name" : "M2",
						"unit" : {
						  "id" : 11,
						  "name" : ""
						},
						"value" : 0.53
				  	}
				  ]
				} 
			]`,
		},
		"Checkup with results, markers and units - return response with results": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),

			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES 
				  	(99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'owner');`,
				`INSERT INTO checkups (id, profile_id, lab_id, status, date) 
					VALUES (10, 3, 0, 'TEST', '2025-07-18');`,
				`INSERT INTO checkup_results (checkup_id, marker_id, unit_id, value) VALUES 
			  		(10, 55, 33, 100),
			  		(10, 77, 11, 0.53);`,
				`INSERT INTO markers (id, name, ref_range_min, ref_range_max) VALUES 
			 		(55, 'M1', 10, 50),
			 		(77, 'M2', 10, 50);`,
				`INSERT INTO units (id, name, unit) VALUES 
				 	(33, 'u1', 'ml'),
				 	(11, 'u2', '%');`,
			},
			expectedCode: 200,
			expectedResp: `[ 
				{
				  "checkup" : {
					"date" : "2025-07-18 00:00:00 +0000 UTC",
					"id" : 10,
					"lab" : {
					  "id" : 0,
					  "name" : ""
					},
					"profile" : {
					  "id" : "00000000-f000-0000-0000-000000000010",
					  "first_name" : "Associated"
					}
				  },
				  "results" : [ 
					{
						"id" : 55,
						"name" : "M1",
						"unit" : {
						  "id" : 33,
						  "name" : "u1"
						},
						"value" : 100
				  	},
					{
						"id" : 77,
						"name" : "M2",
						"unit" : {
						  "id" : 11,
						  "name" : "u2"
						},
						"value" : 0.53
				  	}
				  ]
				} 
			]`,
		},
		"Search query with results - return response with results": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			queryString: "search=m1",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),

			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES 
				  	(99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'owner');`,
				`INSERT INTO checkups (id, profile_id, lab_id, status, date) 
					VALUES (10, 3, 0, 'TEST', '2025-07-18');`,
				`INSERT INTO checkup_results (checkup_id, marker_id, unit_id, value) VALUES 
			  		(10, 55, 33, 100),
			  		(10, 77, 11, 0.53);`,
				`INSERT INTO markers (id, name, ref_range_min, ref_range_max) VALUES 
			 		(55, 'M1', 10, 50),
			 		(77, 'M2', 10, 50);`,
				`INSERT INTO units (id, name, unit) VALUES 
				 	(33, 'u1', 'ml'),
				 	(11, 'u2', '%');`,
			},
			expectedCode: 200,
			expectedResp: `[ 
				{
				  "checkup" : {
					"date" : "2025-07-18 00:00:00 +0000 UTC",
					"id" : 10,
					"lab" : {
					  "id" : 0,
					  "name" : ""
					},
					"profile" : {
					  "id" : "00000000-f000-0000-0000-000000000010",
					  "first_name" : "Associated"
					}
				  },
				  "results" : [ 
					{
						"id" : 55,
						"name" : "M1",
						"unit" : {
						  "id" : 33,
						  "name" : "u1"
						},
						"value" : 100
				  	}
				  ]
				} 
			]`,
		},
		"Search query without results - return response with empty results": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			queryString: "search=noresults",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),

			initialDBQueries: []string{
				`INSERT INTO users(id, uuid, l_name, password) VALUES 
				  	(99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'owner');`,
				`INSERT INTO checkups (id, profile_id, lab_id, status, date) 
					VALUES (10, 3, 0, 'TEST', '2025-07-18');`,
				`INSERT INTO checkup_results (checkup_id, marker_id, unit_id, value) VALUES 
			  		(10, 55, 33, 100),
			  		(10, 77, 11, 0.53);`,
				`INSERT INTO markers (id, name, ref_range_min, ref_range_max) VALUES 
			 		(55, 'M1', 10, 50),
			 		(77, 'M2', 10, 50);`,
				`INSERT INTO units (id, name, unit) VALUES 
				 	(33, 'u1', 'ml'),
				 	(11, 'u2', '%');`,
			},
			expectedCode: 200,
			expectedResp: `[]`,
		},
		//TODO: filters
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
			req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8099/api/v1/account/profiles/%s/checkups?%s", tc.profileUUID, tc.queryString), nil)
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
