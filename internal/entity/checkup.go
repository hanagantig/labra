package entity

import "time"

const (
	UnverifiedCheckupStatus = "unverified"
	VerifiedCheckupStatus   = "verified"
)

type CheckupStatus string

type Checkups []Checkup
type Checkup struct {
	ID             int
	Title          string
	UploadedFileID int
	Profile        Profile
	Lab            Lab
	Status         CheckupStatus
	Date           time.Time
	Comment        string
}

func (c Checkups) IDs() []int {
	ids := make([]int, len(c))
	for i, checkup := range c {
		ids[i] = checkup.ID
	}

	return ids
}

type CheckupResults struct {
	Checkup Checkup
	Results MarkerResults
}

//func NeChecupsWithResilts(checkups Checkups, results Results)
