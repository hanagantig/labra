package checkup

import (
	"database/sql"
	"github.com/google/uuid"
	"labra/internal/entity"
	"time"
)

type Checkup struct {
	ID             int            `db:"id"`
	Title          string         `db:"title"`
	ProfileID      int            `db:"profile_id"`
	ProfileUUID    string         `db:"profile_uuid"`
	PatientFName   sql.NullString `db:"f_name"`
	PatientLName   sql.NullString `db:"l_name"`
	LabID          int            `db:"lab_id"`
	LabName        sql.NullString `db:"lab_name"`
	Status         string         `db:"status"`
	UploadedFileID sql.NullInt64  `db:"uploaded_file_id"`
	Date           time.Time      `db:"date"`
	Comment        string         `db:"comment"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Checkups []Checkup

func NewCheckupFromEntity(e entity.Checkup) Checkup {
	return Checkup{
		ID:             e.ID,
		Title:          e.Title,
		ProfileID:      e.Profile.ID,
		LabID:          e.Lab.ID,
		Date:           e.Date,
		Status:         string(e.Status),
		Comment:        "",
		UploadedFileID: sql.NullInt64{Valid: e.UploadedFileID != 0, Int64: int64(e.UploadedFileID)},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func (c Checkup) BuildEntity() entity.Checkup {
	pUUID, _ := uuid.Parse(c.ProfileUUID)

	return entity.Checkup{
		ID:    c.ID,
		Title: c.Title,
		Profile: entity.Profile{
			ID:    c.ProfileID,
			Uuid:  pUUID,
			FName: c.PatientFName.String,
			LName: c.PatientLName.String,
		},
		Date: c.Date,
		Lab: entity.Lab{
			ID:   c.LabID,
			Name: c.LabName.String,
		},
		UploadedFileID: int(c.UploadedFileID.Int64),
		Status:         entity.CheckupStatus(c.Status),
	}
}

func (c Checkups) BuildEntity() []entity.Checkup {
	if len(c) == 0 {
		return nil
	}

	result := make([]entity.Checkup, len(c))
	for i, ch := range c {
		result[i] = ch.BuildEntity()
	}

	return result
}
