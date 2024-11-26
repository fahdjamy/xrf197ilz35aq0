package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"xrf197ilz35aq0/internal"
)

func TestNewSettings(t *testing.T) {
	t.Run("invalid params", func(t *testing.T) {
		var encryptionKey = internal.RandomBytes(32)
		userFP := "user-fingerprint"
		settings := NewSettings(true, time.Hour, userFP, string(encryptionKey))
		assert.NotNil(t, settings)
	})
}
