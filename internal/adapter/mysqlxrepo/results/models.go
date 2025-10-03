package results

import (
	"database/sql"
	"labra/internal/entity"
	"time"
)

type Results []Result

type Result struct {
	ID              int            `db:"id"`
	CheckupID       int            `db:"checkup_id"`
	MarkerID        int            `db:"marker_id"`
	UndefinedMarker sql.NullString `db:"undefined_marker"`
	UnitID          int            `db:"unit_id"`
	UndefinedUnit   sql.NullString `db:"undefined_unit"`
	Value           float64        `db:"value"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
}

type MarkerResult struct {
	ID              int            `db:"id"`
	CheckupID       int            `db:"checkup_id"`
	MarkerID        int            `db:"marker_id"`
	MarkerName      sql.NullString `db:"marker_name"`
	UndefinedMarker sql.NullString `db:"undefined_marker"`
	MarkerValue     float64        `db:"value"`
	UnitID          int            `db:"unit_id"`
	UnitName        sql.NullString `db:"unit_name"`
	UndefinedUnit   sql.NullString `db:"undefined_unit"`
	PrimaryColor    sql.NullInt64  `db:"primary_color"`
	CreatedAt       time.Time      `db:"created_at"`
}

type MarkerResults []MarkerResult

func NewResultsFromEntity(checkupID int, e entity.MarkerResults) Results {
	res := make(Results, 0, len(e))

	for _, m := range e {
		item := Result{
			CheckupID: checkupID,
			MarkerID:  m.ID,
			UnitID:    m.Unit.ID,
			Value:     m.Value,
		}
		if m.ID == 0 {
			item.UndefinedMarker = sql.NullString{Valid: true, String: m.Marker.Name}
		}

		if m.Unit.ID == 0 {
			item.UndefinedUnit = sql.NullString{Valid: true, String: m.Unit.Name}
		}

		res = append(res, item)
	}

	return res
}

// BuildEntity converts a MarkerResult DB model to the new entity.MarkerResult structure.
func (m MarkerResult) BuildEntity() entity.MarkerResult {
	return entity.MarkerResult{
		ID:        m.ID,
		CheckupID: m.CheckupID,
		Marker: entity.Marker{
			ID:             m.MarkerID,
			Name:           m.MarkerName.String,
			ReferenceRange: entity.Range{}, // Fill if available
			PrimaryColor:   int(m.PrimaryColor.Int64),
		},
		UnrecognizedName: m.UndefinedMarker.String,
		Value:            m.MarkerValue,
		Unit: entity.Unit{
			ID:               m.UnitID,
			Name:             m.UnitName.String,
			UnrecognizedName: m.UndefinedUnit.String,
		},
		CreatedAt: m.CreatedAt,
	}
}

func (m MarkerResults) BuildEntity() entity.MarkerResults {
	res := make(entity.MarkerResults, 0, len(m))
	for _, item := range m {
		res = append(res, item.BuildEntity())
	}
	return res
}
