package unit

import "labra/internal/entity"

type Units []Unit
type Unit struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Unit        string `db:"unit"`
	Description string `db:"description"`
}

func (u Unit) ToEntity() entity.Unit {
	return entity.Unit{
		ID:          u.ID,
		Name:        u.Name,
		FullName:    u.Unit,
		Description: u.Description,
	}
}

func (u Units) ToEntity() entity.Units {
	var units entity.Units
	for _, unit := range u {
		units = append(units, unit.ToEntity())
	}
	return units
}
