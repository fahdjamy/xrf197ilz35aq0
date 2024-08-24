package model

import (
	"time"
	"xrf197ilz35aq0/random"
)

const fingerPrintLength = 31

type User struct {
	id          int64
	FirstName   string
	LastName    string
	email       string
	fingerPrint string
	createdAt   time.Time
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

	return &User{
		id:          id,
		createdAt:   now,
		email:       email,
		LastName:    lastName,
		FirstName:   firstName,
		fingerPrint: uniqueStr,
	}, nil
}
