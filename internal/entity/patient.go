package entity

import (
	"github.com/google/uuid"
	"time"
)

type Profiles []Profile

type Profile struct {
	ID               int
	Uuid             uuid.UUID
	UserID           int
	CreatorUserID    int
	FName            string
	LName            string
	Gender           string
	DateOfBirth      time.Time
	LinkedUserID     int
	LinkedUserAccess string
	Contacts         []Contact
}

func NewProfileFromUser(user User) Profile {
	return Profile{
		UserID:        user.ID,
		CreatorUserID: user.ID,
		FName:         user.FName,
		LName:         user.LName,
		Gender:        user.Gender,
		DateOfBirth:   user.BirthDate,
	}
}

func (p Profile) FullName() string {
	return p.FName + " " + p.LName
}

func (p Profile) IsBoundToUser(userID int) bool {
	return p.LinkedUserID == userID
}

func (p Profile) IsOwnedByUser(userID int) bool {
	ownAsCreator := p.UserID == 0 && p.CreatorUserID == 0
	ownAsLinkedOwner := p.LinkedUserID == userID && p.LinkedUserAccess == "owner"

	return ownAsCreator || ownAsLinkedOwner
}

func (p Profiles) GetIDs() []int {
	ids := make([]int, 0, len(p))
	for _, patient := range p {
		ids = append(ids, patient.ID)
	}

	return ids
}

func (p Profiles) HasProfileUUID(profileUUID uuid.UUID) bool {
	return p.ProfileByUUID(profileUUID).ID > 0
}

func (p Profiles) ProfileByUUID(profileUUID uuid.UUID) Profile {
	for _, profile := range p {
		if profile.Uuid == profileUUID {
			return profile
		}
	}

	return Profile{}
}

func (p Profile) HasAssociatedUser() bool {
	return p.UserID != 0
}
