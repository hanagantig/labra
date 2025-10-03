package cmd

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"labra/internal/entity"
	"reflect"
	"strings"
	"testing"
	"time"
)

type expectedUser struct {
	ID       int
	Uuid     string
	Password string
}

type expectedProfile struct {
	ID            int
	Uuid          string
	UserID        int    `db:"user_id"`
	FName         string `db:"f_name"`
	Gender        string `db:"gender"`
	CreatorUserID int    `db:"creator_user_id"`
}

type expectedUserProfile struct {
	UserID      int    `db:"user_id"`
	ProfileID   int    `db:"profile_id"`
	AccessLevel string `db:"access_level"`
}

type expectedContact struct {
	ID    int
	Type  string
	Value string
}

type expectedLinkedContact struct {
	ID         int
	ContactID  int    `db:"contact_id"`
	EntityType string `db:"entity_type"`
	EntityID   string `db:"entity_id"`
	VerifiedAt string `db:"verified_at"`
}

type expectedCode struct {
	UserID     int    `db:"user_id"`
	ObjectType string `db:"object_type"`
	ObjectID   string `db:"object_id"`
	Code       int    `db:"code"`
}

type expectedUploadedFile struct {
	UserID       int    `db:"user_id"`
	ProfileID    int    `db:"profile_id"`
	FileID       string `db:"file_id"`
	PipelineID   string `db:"pipeline_id"`
	Fingerprint  string `db:"fingerprint"`
	FileType     string `db:"file_type"`
	Source       string `db:"source"`
	Status       string `db:"status"`
	AttemptsLeft int    `db:"attempts_left"`
	Details      string `db:"details"`
}

type expectedCheckup struct {
	ProfileID      int    `db:"profile_id"`
	LabID          int    `db:"lab_id"`
	UploadedFileID string `db:"uploaded_file_id"`
	Status         string `db:"status"`
}

type expectedCheckupResult struct {
	CheckupID       int    `db:"checkup_id"`
	MarkerID        int    `db:"marker_id"`
	UndefinedMarker string `db:"undefined_marker"`
	UnitID          int    `db:"unit_id"`
	UndefinedUnit   string `db:"undefined_unit"`
	Value           string `db:"value"`
}

func newUserToken(uuid uuid.UUID, secret string) string {
	jwt, _ := entity.NewUserJWT(entity.Session{
		UserUUID:  uuid,
		SessionID: uuid,
	}, secret, 15*time.Minute)

	return jwt.String()
}

func assertJSONEqual(t *testing.T, expectedJSON, actualJSON string) {
	var expected, actual interface{}

	dec1 := json.NewDecoder(strings.NewReader(expectedJSON))
	dec1.UseNumber()
	err := dec1.Decode(&expected)
	if err != nil {
		require.Equal(t, expectedJSON, actualJSON)

		return
	}
	require.NoError(t, err, "failed to unmarshal expected JSON")

	dec2 := json.NewDecoder(strings.NewReader(actualJSON))
	dec2.UseNumber()
	err = dec2.Decode(&actual)
	require.NoError(t, err, "failed to unmarshal actual JSON")

	if !reflect.DeepEqual(expected, actual) {
		expB, _ := json.MarshalIndent(expected, "", "  ")
		actB, _ := json.MarshalIndent(actual, "", "  ")
		t.Errorf("response mismatch\nExpected:\n%s\n\nGot:\n%s", expB, actB)
	}

	require.Equal(t, expected, actual, "JSON mismatch")
}

func getAllTableNames(db *sqlx.DB, dbName string) ([]string, error) {
	var tables []string
	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_type = 'BASE TABLE'`
	err := db.Select(&tables, query, dbName)
	return tables, err
}

func truncateAllTables(db *sqlx.DB, dbName string) {
	tables, err := getAllTableNames(db, dbName)
	if err != nil {
		panic(err)
	}

	dbx, err := db.Beginx()
	if err != nil {
		panic(err)
	}

	// Disable foreign key checks
	dbx.MustExec("SET FOREIGN_KEY_CHECKS = 0")

	for _, table := range tables {
		dbx.MustExec("TRUNCATE TABLE " + table)
	}

	err = dbx.Commit()
	if err != nil {
		panic(err)
	}

	// Re-enable foreign key checks
	db.MustExec("SET FOREIGN_KEY_CHECKS = 1")
}
