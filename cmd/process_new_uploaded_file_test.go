package cmd

import (
	"context"
	"errors"
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

func TestCron_ProcessNewUploadedFile(t *testing.T) {
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

		expectedUploadedFiles []expectedUploadedFile
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
		"uploaded file in recognizing status - do nothing": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'recognizing', 'fingerprint1', 5, now(), now());`,
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
					Status:       "recognizing",
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
		"new uploaded file with 0 attempts left - do nothing": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, created_at, updated_at)
					VALUES 
					    (99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'new', 'fingerprint1', now(), now());`,
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
					Status:      "new",
				},
			},
		},
		"new uploaded file fail to check for existing recognition - return error": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'new', 'fingerprint1', 1, now(), now());`,
			},
			getStandardizationCode: http.StatusInternalServerError,
			getStandardizationResp: []byte{},

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
					AttemptsLeft: 1,
				},
			},

			expectedError: errors.New("unexpected status code: 500"),
		},
		"new uploaded file start recognition request fail - return error": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'new', 'fingerprint1', 1, now(), now());`,
			},
			getStandardizationCode: http.StatusNotFound,
			getStandardizationResp: []byte{},

			postStandardizationCode: http.StatusInternalServerError,
			postStandardizationResp: []byte{},

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "file_id1",
					PipelineID:   "",
					Fingerprint:  "fingerprint1",
					FileType:     "file_type",
					Source:       "s",
					Status:       "new",
					AttemptsLeft: 0,
					Details:      "failed to recognize document: ",
				},
			},

			expectedError: errors.New("failed to recognize document: "),
		},
		"new uploaded - run recognition returns empty pipeline id - return error": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'new', 'fingerprint1', 1, now(), now());`,
			},
			getStandardizationCode: http.StatusNotFound,
			getStandardizationResp: []byte{},

			postStandardizationCode: http.StatusOK,
			postStandardizationResp: []byte(`{"standardizationIds":[]}`),

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "file_id1",
					PipelineID:   "",
					Fingerprint:  "fingerprint1",
					FileType:     "file_type",
					Source:       "s",
					Status:       "new",
					AttemptsLeft: 0,
					Details:      "no standardizations found: item not found: client error",
				},
			},

			expectedError: fmt.Errorf("no standardizations found: %w", apperror.ErrNotFound),
		},
		"new uploaded - file already has runed recognition - save with existing pipeline id": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'new', 'fingerprint1', 1, now(), now());`,
			},
			getStandardizationCode: http.StatusOK,
			getStandardizationResp: []byte(`[{"standardizationId":"testExisting"}]`),

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "file_id1",
					PipelineID:   "testExisting",
					Fingerprint:  "fingerprint1",
					FileType:     "file_type",
					Source:       "s",
					Status:       "recognizing",
					AttemptsLeft: 0,
				},
			},
		},
		"new uploaded and start new recognition - save with new pipeline id": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'new', 'fingerprint1', 1, now(), now());`,
			},
			getStandardizationCode: http.StatusNotFound,

			postStandardizationCode: http.StatusOK,
			postStandardizationResp: []byte(`{"standardizationIds":["testNewPipeline"]}`),

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "file_id1",
					PipelineID:   "testNewPipeline",
					Fingerprint:  "fingerprint1",
					FileType:     "file_type",
					Source:       "s",
					Status:       "recognizing",
					AttemptsLeft: 0,
				},
			},
		},
		"multiple new uploaded files, start new recognitions - save with new pipeline id": {
			initialDBQueries: []string{
				`INSERT INTO uploaded_files (user_id, profile_id, pipeline_id, file_id, file_type, source, status, fingerprint, attempts_left, created_at, updated_at)
					VALUES 
					    (99, 3, 'pipeline_id1', 'file_id1', 'file_type', 's', 'new', 'fingerprint1', 1, now(), now()),
					    (99, 3, 'pipeline_id1', 'file_id2', 'file_type', 's', 'new', 'fingerprint2', 1, now(), now());`,
			},
			getStandardizationCode: http.StatusNotFound,

			postStandardizationCode: http.StatusOK,
			postStandardizationResp: []byte(`{"standardizationIds":["testNewPipeline"]}`),

			expectedUploadedFiles: []expectedUploadedFile{
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "file_id1",
					PipelineID:   "testNewPipeline",
					Fingerprint:  "fingerprint1",
					FileType:     "file_type",
					Source:       "s",
					Status:       "recognizing",
					AttemptsLeft: 0,
				},
				{
					UserID:       99,
					ProfileID:    3,
					FileID:       "file_id2",
					PipelineID:   "testNewPipeline",
					Fingerprint:  "fingerprint2",
					FileType:     "file_type",
					Source:       "s",
					Status:       "recognizing",
					AttemptsLeft: 0,
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
				if r.Method == "GET" && r.URL.Path == "/standardizations" {
					w.WriteHeader(tc.getStandardizationCode)
					w.Write(tc.getStandardizationResp)
				}

				if r.Method == "POST" && r.URL.Path == "/standardize/batch" {
					w.WriteHeader(tc.postStandardizationCode)
					w.Write(tc.postStandardizationResp)
				}
			}))

			l, err := net.Listen("tcp", "localhost:5577")
			assert.NoError(t, err)

			err = ts.Listener.Close()
			assert.NoError(t, err)

			ts.Listener = l

			ts.Start()
			defer ts.Close()

			err = a.UseCases().ProcessNewUploadedFiles(context.Background())
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
		})
	}
}
