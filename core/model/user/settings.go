package user

import (
	"time"
	"xrf197ilz35aq0/core"
)

// Settings defines the fields that indicate how a user wants their data to be handled/stored or
// presented to the outside world
type Settings struct {
	// if turned on, user encryption key should be rotated
	RotateEncryptionKey bool
	// should be specified in months and should only be set to run during less peak hours
	encryptAfter    time.Duration
	userFingerprint string
	lastModified    time.Time
}

func NewSettings(rotateEncKey bool, encryptAfter time.Duration, userFP string) (*Settings, error) {
	now := time.Now()
	futureTime := time.Now().Add(encryptAfter) // Convert encryptAfter to time.Time
	if now.After(futureTime) {
		return nil, core.InvalidRequest{
			Message: "encryptAfter can not be in the past",
		}
	}

	if isBeforeMonths(encryptAfter, 3) {
		return nil, core.InvalidRequest{
			Message: "encryptAfter should at least be 3 months from now",
		}
	}

	return &Settings{
		lastModified:        now,
		userFingerprint:     userFP,
		encryptAfter:        encryptAfter,
		RotateEncryptionKey: rotateEncKey,
	}, nil
}

func isBeforeMonths(date time.Duration, month int) bool {
	return !iAfterMonths(date, month)
}

func iAfterMonths(date time.Duration, month int) bool {
	monthsFromNow := time.Now().AddDate(0, month, 0)
	dateToCheck := time.Now().Add(date)
	isAfter := dateToCheck.After(monthsFromNow)
	return isAfter
}
