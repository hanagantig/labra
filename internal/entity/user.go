package entity

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type UserEmail string
type UserPassword string

type EmailOrPhone string

func (e EmailOrPhone) String() string {
	return string(e)
}

func (e EmailOrPhone) IsEmail() bool {
	return strings.Contains(e.String(), "@")
}

func (e EmailOrPhone) IsPhone() bool {
	return !e.IsEmail() && strings.Contains(e.String(), "+")
}

func (u UserEmail) String() string {
	return string(u)
}

func (u UserPassword) String() string {
	return string(u)
}

func (u UserPassword) Hashed() (UserPassword, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return UserPassword(hashed), nil
}

func (u UserPassword) IsMatches(psw UserPassword) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.String()), []byte(psw.String())) == nil
}

type User struct {
	ID        int
	Uuid      uuid.UUID
	FName     string
	LName     string
	Gender    string
	BirthDate time.Time
	Password  UserPassword
	Contacts  Contacts
}

type UserWithProfiles struct {
	User
	Profiles Profiles
}

func (u *User) String() string {
	return fmt.Sprintf("%v %v", u.FName, u.LName)
}

func (u *User) IDString() string {
	return fmt.Sprintf("%v", u.ID)
}
