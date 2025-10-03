package cmd

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hanagantig/gracy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"labra/internal/app"
	"labra/internal/entity"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPI_PostCheckupOCR(t *testing.T) {
	a, err := app.NewApp("../config/app.conf.test.yaml")
	require.NoError(t, err)

	testUserUUID, _ := uuid.Parse("1bba9769-5f72-11f0-93b8-0242ac110002")

	tests := map[string]struct {
		initialDBQueries []string
		authToken        string
		profileUUID      string

		storeDocCode int
		//[]byte(`{"documentId":"test123"}`)
		storeDocResp []byte

		standardizeCode int
		standardizeResp []byte

		fileBytes []byte
		fileName  string

		expectedCode int
		expectedResp string

		expectedUploadedFiles []expectedUploadedFile
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
			fileBytes:    []byte("this is test content of a fake file"),
			fileName:     "test_checkup.jpg",
			expectedCode: 400,
			expectedResp: `{
				"code":400, 
				"error": {
					"message":"unable to parse profile id: invalid UUID length: 4: bad request: client error"
				}
			}`,
		},
		"No user - return 404": {
			profileUUID: uuid.New().String(),
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			fileBytes:   genPNG1x1(),
			fileName:    "test_checkup.png",

			initialDBQueries: []string{
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'owner');`,
			},
			expectedCode: 404,
			expectedResp: `{
				"code":404, 
				"error": {
					"message":"unable to get user profiles: no profiles found for uuid: item not found: client error"
				}
			}`,
		},
		"No user profile - return 404": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			fileBytes:   genPDFTiny(),
			fileName:    "test_checkup.pdf",

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
					"message":"unable to get user profiles: no profiles found for uuid: item not found: client error"
				}
			}`,
		},
		"No profiles - return 404": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			fileBytes:   genJPEGTiny(),
			fileName:    "test_checkup.jpg",

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
					"message":"unable to get user profiles: no profiles found for uuid: item not found: client error"
				}
			}`,
		},
		"invalid file - return 400": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			fileBytes:   randomBytes(49),
			fileName:    "test_checkup.jpg",

			expectedCode: 400,
			expectedResp: `{
				"code":400, 
				"error": {
					"message":"invalid file type: application/octet-stream: bad request: client error"
				}
			}`,
		},
		"file is too big - return 400": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),

			fileBytes:    randomBytes(20 << 20),
			expectedCode: 400,
			expectedResp: `{
				"code":400, 
				"error": {
					"message":"Unable to parse multipart form: multipart: message too large: bad request: client error"
				}
			}`,
		},
		"same file already uploaded for different account - return 400": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			fileBytes:   genPDFTiny(),
			fileName:    "test_checkup.pdf",

			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, created_at, updated_at)
					VALUES (1, 1, 'pipeline_id', 'file_id', 'file_type', 's', 'status', '` + entity.NewUploadedFile(genPDFTiny()).Fingerprint + `', now(), now());`,
				`INSERT INTO users(id, uuid, l_name, password) VALUES 
				  	(99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'owner');`,
			},
			expectedCode: 400,
			expectedResp: `{
				"code":400, 
				"error": {
					"message":"unable to upload file: unable to verify file: file with fingerprint already exists: bad request: client error"
				}
			}`,

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:      1,
					ProfileID:   1,
					FileID:      "file_id",
					PipelineID:  "pipeline_id",
					Fingerprint: entity.NewUploadedFile(genPDFTiny()).Fingerprint,
					FileType:    "file_type",
					Source:      "s",
					Status:      "status",
				},
			},
		},
		"same file already uploaded for same account and profile - return 400": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			fileBytes:   genPDFTiny(),
			fileName:    "test_checkup.pdf",

			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, created_at, updated_at)
					VALUES (99, 3, 'pipeline_id', 'file_id', 'file_type', 's', 'status', '` + entity.NewUploadedFile(genPDFTiny()).Fingerprint + `', now(), now());`,
				`INSERT INTO users(id, uuid, l_name, password) VALUES 
				  	(99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'owner');`,
			},
			expectedCode: 400,
			expectedResp: `{
				"code":400, 
				"error": {
					"message":"unable to upload file: unable to verify file: file with fingerprint already exists: bad request: client error"
				}
			}`,

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:      99,
					ProfileID:   3,
					FileID:      "file_id",
					PipelineID:  "pipeline_id",
					Fingerprint: entity.NewUploadedFile(genPDFTiny()).Fingerprint,
					FileType:    "file_type",
					Source:      "s",
					Status:      "status",
				},
			},
		},
		"same file already uploaded for the account but different profile - return 400": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			fileBytes:   genPDFTiny(),
			fileName:    "test_checkup.pdf",

			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, created_at, updated_at)
					VALUES (99, 5, 'pipeline_id', 'file_id', 'file_type', 's', 'status', '` + entity.NewUploadedFile(genPDFTiny()).Fingerprint + `', now(), now());`,
				`INSERT INTO users(id, uuid, l_name, password) VALUES 
				  	(99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'owner');`,
			},
			expectedCode: 400,
			expectedResp: `{
				"code":400, 
				"error": {
					"message":"unable to upload file: unable to verify file: file with fingerprint already exists: bad request: client error"
				}
			}`,

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:      99,
					ProfileID:   5,
					FileID:      "file_id",
					PipelineID:  "pipeline_id",
					Fingerprint: entity.NewUploadedFile(genPDFTiny()).Fingerprint,
					FileType:    "file_type",
					Source:      "s",
					Status:      "status",
				},
			},
		},
		"weekly upload limit exceeded for user - return 503": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			fileBytes:   genPDFTiny(),
			fileName:    "test_checkup.pdf",

			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, created_at, updated_at)
					VALUES 
					    (99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'status', 'fingerprint1', now(), now()),
					    (99, 3, 'pipeline_id2', 'file_id2', 'file_type', 's', 'status', 'fingerprint2', now(), now()),
					    (99, 3, 'pipeline_id3', 'file_id3', 'file_type', 's', 'status', 'fingerprint3', now(), now()),
					    (99, 3, 'pipeline_id4', 'file_id4', 'file_type', 's', 'status', 'fingerprint4', now(), now()),
					    (99, 3, 'pipeline_id5', 'file_id5', 'file_type', 's', 'status', 'fingerprint5', now(), now());`,
				`INSERT INTO users(id, uuid, l_name, password) VALUES 
				  	(99, '1bba9769-5f72-11f0-93b8-0242ac110002','f name', 'pass')`,
				`INSERT INTO profiles (id, uuid, user_id, creator_user_id, f_name, l_name) VALUES 
				  	(3, '00000000-f000-0000-0000-000000000010', 99, 99, 'Associated', 'Profile'),
				  	(5, '00000000-f000-0000-0000-000000000050', null, 99, 'Added', 'Profile');`,
				`INSERT INTO user_profiles (id, user_id, profile_id, access_level) VALUES 
				  	(11,99, 3, 'owner'),
				  	(15,99, 5, 'owner');`,
			},
			expectedCode: 429,
			expectedResp: `{
				"code":429, 
				"error": {
					"message":"unable to upload file: unable to verify file: uploaded files weekly limit exceeded: Too Many Requests: client error: bad request: client error"
				}
			}`,

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:      99,
					ProfileID:   3,
					FileID:      "file_id1",
					PipelineID:  "pipeline_id1",
					Fingerprint: "fingerprint1",
					FileType:    "file_type",
					Source:      "s",
					Status:      "status",
				},
				{
					UserID:      99,
					ProfileID:   3,
					FileID:      "file_id2",
					PipelineID:  "pipeline_id2",
					Fingerprint: "fingerprint2",
					FileType:    "file_type",
					Source:      "s",
					Status:      "status",
				},
				{
					UserID:      99,
					ProfileID:   3,
					FileID:      "file_id3",
					PipelineID:  "pipeline_id3",
					Fingerprint: "fingerprint3",
					FileType:    "file_type",
					Source:      "s",
					Status:      "status",
				},
				{
					UserID:      99,
					ProfileID:   3,
					FileID:      "file_id4",
					PipelineID:  "pipeline_id4",
					Fingerprint: "fingerprint4",
					FileType:    "file_type",
					Source:      "s",
					Status:      "status",
				},
				{
					UserID:      99,
					ProfileID:   3,
					FileID:      "file_id5",
					PipelineID:  "pipeline_id5",
					Fingerprint: "fingerprint5",
					FileType:    "file_type",
					Source:      "s",
					Status:      "status",
				},
			},
		},
		"failed to upload file to storage - return 500": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			fileBytes:   genPDFTiny(),
			fileName:    "test_checkup.pdf",

			storeDocCode: http.StatusBadRequest,
			storeDocResp: []byte(`{"detail":"fail"}`),

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
			expectedCode: 500,
			expectedResp: `{
				"code":500, 
				"error": {
					"message":"unable to upload file: unable to store file: unexpected status code: 400: fail"
				}
			}`,
		},
		"successfully upload pdf file - return 201": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			fileBytes:   genPDFTiny(),
			fileName:    "test_checkup.pdf",

			storeDocCode: http.StatusOK,
			storeDocResp: []byte(`{"documentId":"testFilePDF"}`),

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
			expectedCode: 201,
			expectedResp: `{}`,

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "testFilePDF",
					Fingerprint:  entity.NewUploadedFile(genPDFTiny()).Fingerprint,
					Source:       "docupanda",
					Status:       "new",
					AttemptsLeft: 4,
				},
			},
		},
		"successfully upload jpg file - return 201": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			fileBytes:   genJPEGTiny(),
			fileName:    "test_checkup.jpg",

			storeDocCode: http.StatusOK,
			storeDocResp: []byte(`{"documentId":"testFileJPG"}`),

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
			expectedCode: 201,
			expectedResp: `{}`,

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "testFileJPG",
					Fingerprint:  entity.NewUploadedFile(genJPEGTiny()).Fingerprint,
					Source:       "docupanda",
					Status:       "new",
					AttemptsLeft: 4,
				},
			},
		},
		"successfully upload png file - return 201": {
			profileUUID: "00000000-f000-0000-0000-000000000010",
			authToken:   newUserToken(testUserUUID, a.GetConfig().Auth.AccessTokenSecret),
			fileBytes:   genPDFTiny(),
			fileName:    "test_checkup.png",

			storeDocCode: http.StatusOK,
			storeDocResp: []byte(`{"documentId":"testFilePNG"}`),

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
			expectedCode: 201,
			expectedResp: `{}`,

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "testFilePNG",
					Fingerprint:  entity.NewUploadedFile(genPDFTiny()).Fingerprint,
					Source:       "docupanda",
					Status:       "new",
					AttemptsLeft: 4,
				},
			},
		},

		//"TODO: failed to insert uploaded file data to db - return 500": {},
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

			reqBody := &bytes.Buffer{}
			writer := multipart.NewWriter(reqBody)

			// Create the form file part
			part, err := writer.CreateFormFile("report", tc.fileName)
			if err != nil {
				require.NoError(t, err)
			}

			// Write byte content to form field
			_, err = io.Copy(part, bytes.NewReader(tc.fileBytes))
			if err != nil {
				require.NoError(t, err)
			}

			// Finalize the multipart form
			err = writer.Close()
			if err != nil {
				require.NoError(t, err)
			}

			ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "POST" && r.URL.Path == "/document" {
					w.WriteHeader(tc.storeDocCode)
					w.Write(tc.storeDocResp)
				}

				if r.Method == "POST" && r.URL.Path == "/standardize/batch" {
					w.WriteHeader(tc.storeDocCode)
					w.Write(tc.storeDocResp)
				}
			}))
			defer ts.Close()

			l, err := net.Listen("tcp", "localhost:5577")
			assert.NoError(t, err)

			err = ts.Listener.Close()
			assert.NoError(t, err)

			ts.Listener = l

			ts.Start()

			//httpmock.Activate()
			//defer httpmock.DeactivateAndReset()

			httpClient := &http.Client{}
			req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8099/api/v1/account/profiles/%s/checkups/ocr", tc.profileUUID), reqBody)
			require.NoError(t, err)

			req.Header.Set("Authorization", "Bearer "+tc.authToken)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			res, err := httpClient.Do(req)

			require.NoError(t, err)

			body, _ := io.ReadAll(res.Body)
			defer res.Body.Close()

			require.Equal(t, tc.expectedCode, res.StatusCode, string(body))
			assertJSONEqual(t, tc.expectedResp, string(body))

			var uploadedFiles []expectedUploadedFile
			err = db.Select(&uploadedFiles, `SELECT user_id, profile_id, file_id, pipeline_id, fingerprint, file_type, source, status, attempts_left, coalesce(details,"") as details FROM uploaded_files`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedUploadedFiles, uploadedFiles)

			//var userProfiles []expectedUserProfile
			//err = db.Select(&userProfiles, `SELECT user_id, profile_id, access_level FROM user_profiles;`)
			//require.NoError(t, err)
			//require.Equal(t, tc.expectedUserProfiles, userProfiles)
		})
	}
}

func genPNG1x1() []byte {
	// Minimal valid 1x1 PNG (transparent). Magic + IHDR + IDAT + IEND.
	// (Precomputed tiny PNG bytes)
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0x60, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x01, 0xE5, 0x27, 0xD4, 0xA2, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		0x42, 0x60, 0x82,
	}
}

func genJPEGTiny() []byte {
	// Minimal JFIF JPEG header + EOI. Enough for content sniffers.
	return []byte{
		0xFF, 0xD8, // SOI
		0xFF, 0xE0, 0x00, 0x10, // APP0 marker, len
		'J', 'F', 'I', 'F', 0x00, 0x01, 0x01, 0x00,
		0x00, 0x01, 0x00, 0x01, 0x00, 0x00,
		0xFF, 0xDB, 0x00, 0x43, 0x00, // DQT (truncated table ok for sniff)
		// pad some bytes
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0xFF, 0xD9, // EOI
	}
}

func genGIF1x1() []byte {
	// Minimal GIF87a 1x1 transparent
	return []byte("GIF87a\x01\x00\x01\x00\x80\x01\x00\x00\x00\x00\xFF\xFF\xFF,\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02D\x01\x00;")
}

func genPDFTiny() []byte {
	// Minimal PDF with header, body, xref, trailer, EOF.
	return []byte("%PDF-1.4\n1 0 obj\n<<>>\nendobj\nxref\n0 2\n0000000000 65535 f \n0000000010 00000 n \ntrailer\n<< /Size 2 >>\nstartxref\n38\n%%EOF\n")
}

func genText() []byte {
	return []byte("hello world\n")
}

// randomBytes returns n random bytes (used to test size limits)
func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}
