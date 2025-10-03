package entity

import (
	"crypto/rand"
	"math/big"
	"time"
)

const NewUserContactObjectType = "NEW_USER_CONTACT"

const defaultOTPCodeLen = 5

type OTPObjectType string

func (o OTPObjectType) String() string {
	return string(o)
}

type OTPCode struct {
	ID         uint
	UserID     int
	ObjectType OTPObjectType
	ObjectID   string
	Code       string
	ExpiredAt  time.Time
	CreatedAt  time.Time
}

func NewCode(userID int, objType OTPObjectType, objID string, ttl time.Duration) OTPCode {
	code := OTPCode{
		UserID:     userID,
		ObjectType: objType,
		ObjectID:   objID,
		ExpiredAt:  time.Now().Add(ttl),
	}

	code.Code = code.generateCode(defaultOTPCodeLen)

	return code
}

func NewCodeWithValue(userID int, objType OTPObjectType, objID string, codeValue string, ttl time.Duration) OTPCode {
	return OTPCode{
		UserID:     userID,
		ObjectType: objType,
		ObjectID:   objID,
		Code:       codeValue,
		ExpiredAt:  time.Now().Add(ttl),
		CreatedAt:  time.Now(),
	}
}

func (o OTPCode) String() string {
	return o.Code
}

func (o OTPCode) generateCode(n int) string {
	var numbersRunes = []rune("1234567890")

	b := make([]rune, n)
	for i := range b {
		r := int64(0)
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(numbersRunes))))
		// TODO: log error
		if err == nil {
			r = nBig.Int64()
		}
		b[i] = numbersRunes[r]
	}

	return string(b)
}
