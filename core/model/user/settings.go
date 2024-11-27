package user

import (
	"encoding/json"
	"time"
)

// Settings defines the fields that indicate how a user wants their data to be handled/stored or
// presented to the outside world
type Settings struct {
	// if turned on, user encryption key should be rotated
	RotateEncryptionKey bool
	// should be specified in months and should only be set to run during less peak hours
	CreatedAt       time.Time
	encryptAfter    time.Duration
	userFingerprint string
	LastModified    time.Time
	encryptionKey   string
	UserKey         bool
}

func (s *Settings) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, s)
}

func (s *Settings) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Settings) Key() string {
	return s.encryptionKey
}

func NewSettings(rotateEncKey bool, encryptAfter time.Duration, userFP, encryptionKey string) *Settings {
	now := time.Now()

	return &Settings{
		CreatedAt:           now,
		LastModified:        now,
		userFingerprint:     userFP,
		encryptAfter:        encryptAfter,
		RotateEncryptionKey: rotateEncKey,
		encryptionKey:       encryptionKey,
	}
}
