package entity

type Units []Unit
type Unit struct {
	ID               int
	Name             string
	UnrecognizedName string
	FullName         string
	Description      string
}

func (u Unit) String() string {
	return u.Name
}
