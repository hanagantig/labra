package marker

import (
	"database/sql"
	"labra/internal/entity"
)

type Markers []Marker

type Marker struct {
	ID           int            `db:"id"`
	Name         string         `db:"name"`
	RefRangeMax  sql.NullString `db:"ref_range_max"`
	RefRangeMin  sql.NullString `db:"ref_range_min"`
	PrimaryColor int            `db:"primary_color"`
}

func (m Marker) BuildEntity() entity.Marker {
	return entity.Marker{
		ID:   m.ID,
		Name: m.Name,
		ReferenceRange: entity.Range{
			From: m.RefRangeMin.String,
			To:   m.RefRangeMax.String,
		},
		PrimaryColor: m.PrimaryColor,
	}
}

func (m Markers) BuildEntity() entity.Markers {
	result := make(entity.Markers, 0, len(m))
	for _, marker := range m {
		result = append(result, marker.BuildEntity())
	}

	return result
}
