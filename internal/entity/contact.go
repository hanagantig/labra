package entity

import (
	"fmt"
	"github.com/google/uuid"
)

type ContactType string

const (
	ContactTypeEmail               ContactType = "email"
	ContactTypePhone               ContactType = "phone"
	LinkedContactEntityUserType                = "user"
	LinkedContactEntityProfileType             = "profile"
)

func NewContactFromLogin(login EmailOrPhone) Contact {
	contact := Contact{
		Value: login.String(),
		Type:  ContactTypeEmail,
	}

	if login.IsPhone() {
		contact.Type = ContactTypePhone
	}

	return contact
}

type LinkedEntity struct {
	ID         int
	Uuid       uuid.UUID
	ContactID  int
	IsPrimary  bool
	Type       string
	IsVerified bool
}

type Contacts []Contact

type Contact struct {
	ID             int
	Type           ContactType
	Value          string
	LinkedEntities []LinkedEntity
}

func (l LinkedEntity) IsUserType() bool {
	return l.Type == LinkedContactEntityUserType
}

func (l LinkedEntity) IsProfileType() bool {
	return l.Type == LinkedContactEntityProfileType
}

func (c Contact) BelongsToUser() bool {
	for _, l := range c.LinkedEntities {
		if l.IsUserType() {
			return true
		}
	}

	return false
}

func (c Contact) VerifiedByUser() bool {
	for _, l := range c.LinkedEntities {
		if l.IsUserType() {
			return l.IsVerified
		}
	}

	return false
}

func (c Contact) VerifiedByPatient() bool {
	for _, l := range c.LinkedEntities {
		if l.IsUserType() {
			return l.IsVerified
		}
	}

	return false
}

func (c Contact) BelongsToProfile() bool {
	for _, l := range c.LinkedEntities {
		if l.IsProfileType() {
			return true
		}
	}

	return false
}

func (c Contacts) ProfileEntity() LinkedEntity {
	for _, contact := range c {
		for _, l := range contact.LinkedEntities {
			if l.IsProfileType() {
				return l
			}
		}
	}

	return LinkedEntity{}
}

func (c Contact) LinkedProfile() LinkedEntity {
	for _, l := range c.LinkedEntities {
		if l.IsProfileType() {
			return l
		}
	}

	return LinkedEntity{}
}

func (c Contact) LinkedUser() LinkedEntity {
	for _, l := range c.LinkedEntities {
		if l.IsUserType() {
			return l
		}
	}

	return LinkedEntity{}
}

func (c Contacts) MapByEntityID(entityType string) map[int]Contact {
	m := make(map[int]Contact)
	for _, contact := range c {
		for _, l := range contact.LinkedEntities {
			if l.Type == entityType {
				m[l.ID] = contact
			}
		}
	}

	return m
}

func (c Contact) IDString() string {
	return fmt.Sprintf("%v", c.ID)
}

func (c Contact) IsPhone() bool {
	return c.Type == ContactTypePhone
}

func (c Contact) IsEmail() bool {
	return c.Type == ContactTypeEmail
}
