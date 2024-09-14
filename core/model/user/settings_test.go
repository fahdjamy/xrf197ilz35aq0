package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"xrf197ilz35aq0"
)

func TestNewSettings(t *testing.T) {
	t.Run("invalid params", func(t *testing.T) {
		userFP := "user-fingerprint"
		tests := []struct {
			name         string
			encryptAfter time.Duration
		}{
			{
				name:         "should return an error if encryptAfter is in the past",
				encryptAfter: -10 * time.Minute,
			},
			{name: "should return an error if encryptAfter not 3 months from now", encryptAfter: time.Hour},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				settings, err := NewSettings(true, tt.encryptAfter, userFP)
				xrf197ilz35aq0.AssertError(t, err)
				assert.Nil(t, settings)
			})
		}
	})
}
