package model

import (
	"time"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/random"
)

const fingerPrintLength = 31

type User struct {
	id          int64
	FirstName   string
	LastName    string
	fingerPrint string
	createdAt   time.Time
	email       Secret[*xrf.SerializableString]
}

func newUser(firstName string, lastName string, email string) (*User, error) {
	now := time.Now()

	uniqueStr, err := random.TimeBasedString(now.Unix(), fingerPrintLength)
	if err != nil {
		return nil, err
	}
	id, err := random.Int64FromUUID()
	if err != nil {
		return nil, err
	}
	serializableEmail := xrf.SerializableString(email)

	var secretEmailPtr Secret[*xrf.SerializableString]
	secretEmailPtr = *NewSecret(&serializableEmail)

	return &User{
		id:          id,
		createdAt:   now,
		LastName:    lastName,
		FirstName:   firstName,
		fingerPrint: uniqueStr,
		email:       secretEmailPtr,
	}, nil
}
