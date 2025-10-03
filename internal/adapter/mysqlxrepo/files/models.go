package files

import (
	"database/sql"
	"labra/internal/entity"
	"time"
)

type UploadedFiles []UploadedFile
type UploadedFile struct {
	ID          int64          `db:"id"`
	UserID      int            `db:"user_id"`
	ProfileID   int            `db:"profile_id"`
	FileID      string         `db:"file_id"`
	FileType    string         `db:"file_type"`
	PipelineID  string         `db:"pipeline_id"`
	Fingerprint string         `db:"fingerprint"`
	Status      string         `db:"status"`
	Source      string         `db:"source"`
	Details     sql.NullString `db:"details"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

func NewUploadedFileFromEntity(e entity.UploadedFile) UploadedFile {
	return UploadedFile{
		UserID:      e.UserID,
		ProfileID:   e.ProfileID,
		FileID:      e.FileID,
		PipelineID:  e.PipelineID,
		Fingerprint: e.Fingerprint,
		Source:      e.Source,
		Status:      string(e.Status),
	}
}

func (u UploadedFile) buildEntity() entity.UploadedFile {
	return entity.UploadedFile{
		ID:          int(u.ID),
		Fingerprint: u.Fingerprint,
		UserID:      u.UserID,
		ProfileID:   u.ProfileID,
		FileID:      u.FileID,
		FileType:    u.FileType,
		PipelineID:  u.PipelineID,
		Status:      entity.UploadedFileStatus(u.Status),
		Source:      u.Source,
		Details:     u.Details.String,
	}
}

func (u UploadedFiles) buildEntity() []entity.UploadedFile {
	res := make([]entity.UploadedFile, len(u))
	for i, u := range u {
		res[i] = u.buildEntity()
	}

	return res
}
