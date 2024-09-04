package encryption

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	t.Run("generate Key", func(t *testing.T) {
		testCases := []struct {
			name           string
			keySize        int
			expectingError bool
		}{
			{
				name:           "Valid key size",
				keySize:        32,
				expectingError: false,
			},
			{
				name:           "Invalid key size",
				keySize:        -1,
				expectingError: true,
			},
			{
				name:           "Invalid key size",
				keySize:        0,
				expectingError: true,
			},
			{
				name:           "Too big key size",
				keySize:        300,
				expectingError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				key, err := GenerateKey(tc.keySize)
				if (err != nil) != tc.expectingError {
					t.Errorf("GenerateKey() error = %v, wantErr %v", err, tc.expectingError)
				}
				if !tc.expectingError && len(key) != tc.keySize {
					t.Errorf("Incorrect key length. Got: %d, want: %d", len(key), tc.keySize)
				}
			})
		}
	})
}
