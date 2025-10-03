package lab

import "labra/internal/entity"

type Labs []Lab
type Lab struct {
	ID   int    `db:"id"`
	Name string `db:"lab_name"`
}

func (l Lab) BuildEntity() entity.Lab {
	return entity.Lab{
		ID:   l.ID,
		Name: l.Name,
	}
}

func (l Labs) BuildEntity() []entity.Lab {
	result := make([]entity.Lab, 0, len(l))
	for _, lab := range l {
		result = append(result, lab.BuildEntity())
	}

	return result
}
