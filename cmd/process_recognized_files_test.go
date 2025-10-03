package cmd

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"labra/internal/app"
	"labra/internal/apperror"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCron_ProcessRecognizedFiles(t *testing.T) {
	a, err := app.NewApp("../config/app.conf.test.yaml")
	require.NoError(t, err)

	tests := map[string]struct {
		initialDBQueries []string

		getStandardizationCode int
		getStandardizationResp []byte

		postStandardizationCode int
		postStandardizationResp []byte

		fileBytes []byte
		fileName  string

		expectedError error

		expectedUploadedFiles  []expectedUploadedFile
		expectedCheckups       []expectedCheckup
		expectedCheckupResults []expectedCheckupResult
	}{
		"No uploaded files - run and do nothing": {},
		"uploaded file in undefined status - do nothing": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'status', 'fingerprint1', 5, now(), now());`,
			},
			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "file_id1",
					PipelineID:   "pipeline_id1",
					Fingerprint:  "fingerprint1",
					FileType:     "file_type",
					Source:       "s",
					Status:       "status",
					AttemptsLeft: 5,
				},
			},
		},
		"uploaded file in new status - do nothing": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'new', 'fingerprint1', 5, now(), now());`,
			},
			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "file_id1",
					PipelineID:   "pipeline_id1",
					Fingerprint:  "fingerprint1",
					FileType:     "file_type",
					Source:       "s",
					Status:       "new",
					AttemptsLeft: 5,
				},
			},
		},
		"uploaded file in recognized status - do nothing": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'recognized', 'fingerprint1', 5, now(), now());`,
			},
			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "file_id1",
					PipelineID:   "pipeline_id1",
					Fingerprint:  "fingerprint1",
					FileType:     "file_type",
					Source:       "s",
					Status:       "recognized",
					AttemptsLeft: 5,
				},
			},
		},
		"recognizing uploaded file with 0 attempts left - do nothing": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, created_at, updated_at)
					VALUES 
					    (99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'recognizing', 'fingerprint1', now(), now());`,
			},
			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:      99,
					ProfileID:   3,
					FileID:      "file_id1",
					PipelineID:  "pipeline_id1",
					Fingerprint: "fingerprint1",
					FileType:    "file_type",
					Source:      "s",
					Status:      "recognizing",
				},
			},
		},
		"fail to get recognized results - return error": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'recognizing', 'fingerprint1', 5, now(), now());`,
			},

			getStandardizationCode: http.StatusBadRequest,

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "file_id1",
					PipelineID:   "pipeline_id1",
					Fingerprint:  "fingerprint1",
					FileType:     "file_type",
					Source:       "s",
					Status:       "recognizing",
					Details:      "unexpected status code: 400",
					AttemptsLeft: 4,
				},
			},

			expectedError: fmt.Errorf("get pipeline results 1: unexpected status code: 400"),
		},
		"file already has checkup - return error": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (id, user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (111, 99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'recognizing', 'fingerprint1', 5, now(), now());`,
				`INSERT INTO checkups (profile_id, lab_id, status, uploaded_file_id, date, comment, created_at, updated_at) 
					VALUES (3, 0, 'unverified', 111, now(), '', now(), now());`,
			},

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "file_id1",
					PipelineID:   "pipeline_id1",
					Fingerprint:  "fingerprint1",
					FileType:     "file_type",
					Source:       "s",
					Status:       "duplicated",
					Details:      "checkup for file already exists: entity already exists: client error",
					AttemptsLeft: 4,
				},
			},

			expectedCheckups: []expectedCheckup{
				{
					ProfileID:      3,
					LabID:          0,
					UploadedFileID: "111",
					Status:         "unverified",
				},
			},

			expectedError: fmt.Errorf("checkup for file already exists: %w", apperror.ErrDuplicateEntity),
		},
		"recognized file - store result as checkup with undefined units and markers": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (id, user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (111, 99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'recognizing', 'fingerprint1', 5, now(), now());`,
			},

			getStandardizationCode: http.StatusOK,
			getStandardizationResp: []byte(
				`{
  					"standardizationId" : "pipeline_id1",
  					"documentId" : "file_id1",
  					"data" : {
    					"laboratory" : "BION LAB4U",
    					"phoneNumber" : "8 (800) 555-35-90",
    					"website" : "Lab4U.ru",
    					"patient" : {
      						"name" : "Ханагян Мария Мартиновна",
      						"birthDate" : "2003-11-24",
      						"gender" : "Женский"
    					},
						"testDetails" : {
						  "testType" : "Общий анализ крови (CBC/Diff) с лейкоцитарной формулой",
						  "sampleDate" : "2024-01-31",
						  "sampleTime" : "09:30",
						  "deliveryDate" : "2024-02-01",
						  "reportDate" : "2024-02-01",
						  "orderNumber" : "977939198001",
						  "sampleNumber" : "977939198001"
						},
						"analysisResults" : [ 
							{
							  "testName" : "Гемоглобин",
							  "result" : {
								"amount" : 118,
								"unit" : "г/л"
							  },
							  "referenceRange" : {
								"lower" : 120,
								"upper" : 158
							  }
							},
							{
							  "testName" : "Эритроциты",
							  "result" : {
								"amount" : 4.09,
								"unit" : "10/12/л"
							  },
							  "referenceRange" : {
								"lower" : 3.9,
								"upper" : 5.2
							  }
							}
						]
  					},
				  	"schemaId" : "5188c4f3",
					"schemaName" : "Laboratory Test Report V2",
				  	"jobId" : "N0C9XwGH",
					"dataset" : "U01P01",
					"filename" : "mgo_test",
					"timestamp" : "2025-08-15T06:57:15.117000Z",
					"metadata" : null
				}`,
			),

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "file_id1",
					PipelineID:   "pipeline_id1",
					Fingerprint:  "fingerprint1",
					FileType:     "file_type",
					Source:       "s",
					Status:       "recognized",
					AttemptsLeft: 4,
				},
			},

			expectedCheckups: []expectedCheckup{
				{
					ProfileID:      3,
					LabID:          0,
					UploadedFileID: "111",
					Status:         "unverified",
				},
			},

			expectedCheckupResults: []expectedCheckupResult{
				{
					CheckupID:       1,
					MarkerID:        0,
					UndefinedMarker: "Гемоглобин",
					UnitID:          0,
					UndefinedUnit:   "г/л",
					Value:           "118",
				},
				{
					CheckupID:       1,
					MarkerID:        0,
					UndefinedMarker: "Эритроциты",
					UnitID:          0,
					UndefinedUnit:   "10/12/л",
					Value:           "4.09",
				},
			},
		},
		"recognized file - store result as checkup with units and undefined markers": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (id, user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (111, 99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'recognizing', 'fingerprint1', 5, now(), now());`,
				`INSERT INTO units (id, name, unit) VALUES (1, "г/л", "г/л"), (3, "10^12/л", "10^12/л")`,
			},

			getStandardizationCode: http.StatusOK,
			getStandardizationResp: []byte(
				`{
  					"standardizationId" : "pipeline_id1",
  					"documentId" : "file_id1",
  					"data" : {
    					"laboratory" : "BION LAB4U",
    					"phoneNumber" : "8 (800) 555-35-90",
    					"website" : "Lab4U.ru",
    					"patient" : {
      						"name" : "Ханагян Мария Мартиновна",
      						"birthDate" : "2003-11-24",
      						"gender" : "Женский"
    					},
						"testDetails" : {
						  "testType" : "Общий анализ крови (CBC/Diff) с лейкоцитарной формулой",
						  "sampleDate" : "2024-01-31",
						  "sampleTime" : "09:30",
						  "deliveryDate" : "2024-02-01",
						  "reportDate" : "2024-02-01",
						  "orderNumber" : "977939198001",
						  "sampleNumber" : "977939198001"
						},
						"analysisResults" : [ 
							{
							  "testName" : "Гемоглобин",
							  "result" : {
								"amount" : 118,
								"unit" : "г/л"
							  },
							  "referenceRange" : {
								"lower" : 120,
								"upper" : 158
							  }
							},
							{
							  "testName" : "Эритроциты",
							  "result" : {
								"amount" : 4.09,
								"unit" : "10/12/л"
							  },
							  "referenceRange" : {
								"lower" : 3.9,
								"upper" : 5.2
							  }
							}
						]
  					},
				  	"schemaId" : "5188c4f3",
					"schemaName" : "Laboratory Test Report V2",
				  	"jobId" : "N0C9XwGH",
					"dataset" : "U01P01",
					"filename" : "mgo_test",
					"timestamp" : "2025-08-15T06:57:15.117000Z",
					"metadata" : null
				}`,
			),

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "file_id1",
					PipelineID:   "pipeline_id1",
					Fingerprint:  "fingerprint1",
					FileType:     "file_type",
					Source:       "s",
					Status:       "recognized",
					AttemptsLeft: 4,
				},
			},

			expectedCheckups: []expectedCheckup{
				{
					ProfileID:      3,
					LabID:          0,
					UploadedFileID: "111",
					Status:         "unverified",
				},
			},

			expectedCheckupResults: []expectedCheckupResult{
				{
					CheckupID:       1,
					MarkerID:        0,
					UndefinedMarker: "Гемоглобин",
					UnitID:          1,
					Value:           "118",
				},
				{
					CheckupID:       1,
					MarkerID:        0,
					UndefinedMarker: "Эритроциты",
					UnitID:          3,
					Value:           "4.09",
				},
			},
		},
		"recognized file - store result as checkup with units and markers": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (id, user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (111, 99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'recognizing', 'fingerprint1', 5, now(), now());`,
				`INSERT INTO units (id, name, unit) VALUES (1, "г/л", "г/л"), (3, "10^12/л", "10^12/л")`,
				`INSERT INTO markers (id, name) VALUES (1, "Гемоглобин"), (7, "Эритроциты")`,
			},

			getStandardizationCode: http.StatusOK,
			getStandardizationResp: []byte(
				`{
  					"standardizationId" : "pipeline_id1",
  					"documentId" : "file_id1",
  					"data" : {
    					"laboratory" : "BION LAB4U",
    					"phoneNumber" : "8 (800) 555-35-90",
    					"website" : "Lab4U.ru",
    					"patient" : {
      						"name" : "Ханагян Мария Мартиновна",
      						"birthDate" : "2003-11-24",
      						"gender" : "Женский"
    					},
						"testDetails" : {
						  "testType" : "Общий анализ крови (CBC/Diff) с лейкоцитарной формулой",
						  "sampleDate" : "2024-01-31",
						  "sampleTime" : "09:30",
						  "deliveryDate" : "2024-02-01",
						  "reportDate" : "2024-02-01",
						  "orderNumber" : "977939198001",
						  "sampleNumber" : "977939198001"
						},
						"analysisResults" : [ 
							{
							  "testName" : "Гемоглобин",
							  "result" : {
								"amount" : 118,
								"unit" : "г/л"
							  },
							  "referenceRange" : {
								"lower" : 120,
								"upper" : 158
							  }
							},
							{
							  "testName" : "Эритроциты",
							  "result" : {
								"amount" : 4.09,
								"unit" : "10/12/л"
							  },
							  "referenceRange" : {
								"lower" : 3.9,
								"upper" : 5.2
							  }
							},
							{
							  "testName" : "Undefined",
							  "result" : {
								"amount" : 55,
								"unit" : "undefined"
							  }
							}
						]
  					},
				  	"schemaId" : "5188c4f3",
					"schemaName" : "Laboratory Test Report V2",
				  	"jobId" : "N0C9XwGH",
					"dataset" : "U01P01",
					"filename" : "mgo_test",
					"timestamp" : "2025-08-15T06:57:15.117000Z",
					"metadata" : null
				}`,
			),

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "file_id1",
					PipelineID:   "pipeline_id1",
					Fingerprint:  "fingerprint1",
					FileType:     "file_type",
					Source:       "s",
					Status:       "recognized",
					AttemptsLeft: 4,
				},
			},

			expectedCheckups: []expectedCheckup{
				{
					ProfileID:      3,
					LabID:          0,
					UploadedFileID: "111",
					Status:         "unverified",
				},
			},

			expectedCheckupResults: []expectedCheckupResult{
				{
					CheckupID: 1,
					MarkerID:  1,
					UnitID:    1,
					Value:     "118",
				},
				{
					CheckupID: 1,
					MarkerID:  7,
					UnitID:    3,
					Value:     "4.09",
				},
				{
					CheckupID:       1,
					UndefinedMarker: "Undefined",
					UndefinedUnit:   "undefined",
					Value:           "55",
				},
			},
		},
	}

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

			ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" && r.URL.Path == "/standardization/pipeline_id1" {
					w.WriteHeader(tc.getStandardizationCode)
					w.Write(tc.getStandardizationResp)
				}
			}))

			l, err := net.Listen("tcp", "localhost:5577")
			assert.NoError(t, err)

			err = ts.Listener.Close()
			assert.NoError(t, err)

			ts.Listener = l

			ts.Start()
			defer ts.Close()

			err = a.UseCases().ProcessRecognizedFiles(context.Background())
			if tc.expectedError != nil {
				require.NotNil(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			var uploadedFiles []expectedUploadedFile
			err = db.Select(&uploadedFiles, `SELECT user_id, profile_id, file_id, pipeline_id, fingerprint, file_type, source, status, attempts_left, coalesce(details,"") as details FROM uploaded_files`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedUploadedFiles, uploadedFiles)

			var checkups []expectedCheckup
			err = db.Select(&checkups, `SELECT profile_id, lab_id, status, uploaded_file_id FROM checkups`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedCheckups, checkups)

			var checkupResults []expectedCheckupResult
			err = db.Select(&checkupResults, `SELECT checkup_id, marker_id, coalesce(undefined_marker, "") as undefined_marker, unit_id, coalesce(undefined_unit, "") as undefined_unit, value FROM checkup_results`)
			require.NoError(t, err)
			require.Equal(t, tc.expectedCheckupResults, checkupResults)
		})
	}
}
