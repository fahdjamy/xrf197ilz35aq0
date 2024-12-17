package user

import (
	"encoding/json"
	"time"
	xrfErr "xrf197ilz35aq0/internal/error"
)

// Settings defines the fields that indicate how a user wants their data to be handled/stored or
// presented to the outside world
type Settings struct {
	// if turned on, user encryption key should be rotated
	RotateEncryptionKey bool
	// should be specified in months and should only be set to run during less peak hours
	CreatedAt       time.Time
	EncryptAfter    time.Duration
	UserFingerprint string
	LastModified    time.Time
	EncryptionKey   string
	UserKey         bool

	// Argon2 parameters
	Time    uint32 `bson:"argon2Time"`
	Memory  uint32 `bson:"argon2Memory"`
	Threads uint8  `bson:"argon2Threads"`
}

func (s *Settings) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, s)
}

func (s *Settings) MarshalJSON() ([]byte, error) {
	internErr := &xrfErr.Internal{
		Message: "internal error",
		Source:  "model/user/settings#MarshalJSON",
	}

	if s == nil {
		internErr.Message = "settings is nil"
		return nil, internErr
	}

	return json.Marshal(struct {
		RotateEncryptionKey bool          `json:"rotateKey"`
		CreatedAt           time.Time     `json:"createdAt"`
		EncryptAfter        time.Duration `json:"encryptAfter"`
		LastModified        time.Time     `json:"lastModified"`
		EncryptionKey       string        `json:"encryptionKey"`
		UserKey             bool          `json:"userKey"`
	}{
		RotateEncryptionKey: s.RotateEncryptionKey,
		CreatedAt:           s.CreatedAt,
		EncryptAfter:        s.EncryptAfter,
		LastModified:        s.LastModified,
		EncryptionKey:       s.EncryptionKey,
		UserKey:             s.UserKey,
	})
}

func (s *Settings) Key() string {
	return s.EncryptionKey
}

func NewSettings(rotateEncKey bool, encryptAfter time.Duration, userFP, encryptionKey string) *Settings {
	now := time.Now()

	return &Settings{
		CreatedAt:           now,
		LastModified:        now,
		UserFingerprint:     userFP,
		EncryptAfter:        encryptAfter,
		RotateEncryptionKey: rotateEncKey,
		EncryptionKey:       encryptionKey,
	}
}
