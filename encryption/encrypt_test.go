package encryption

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"xrf197ilz35aq0"
)

var (
	plaintext     = []byte("Some sensitive data")
	encryptionKey = xrf197ilz35aq0.RandomBytes(32)
)

func TestEncrypt(t *testing.T) {
	tests := []struct {
		shouldErr bool
		name      string
		plaintext []byte
		key       []byte
	}{
		{
			name:      "Should err if key is less than 32",
			key:       xrf197ilz35aq0.RandomBytes(31),
			shouldErr: true,
			plaintext: plaintext,
		},
		{
			name:      "Should err if data is nil",
			plaintext: nil,
			shouldErr: true,
			key:       encryptionKey,
		},
		{
			name:      "Should encrypt plaintext",
			plaintext: plaintext,
			shouldErr: false,
			key:       encryptionKey,
		},
		{
			name:      "Should err of key is nil",
			plaintext: plaintext,
			shouldErr: true,
			key:       nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encryptedData, err := Encrypt(test.plaintext, test.key)
			if test.shouldErr {
				xrf197ilz35aq0.AssertError(t, err)
				assert.Nil(t, encryptedData)
			} else {
				xrf197ilz35aq0.AssertNoError(t, err)
				assert.NotNil(t, encryptedData)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	encryptedData, err := Encrypt(plaintext, encryptionKey)
	xrf197ilz35aq0.AssertNoError(t, err)

	tests := []struct {
		shouldErr bool
		name      string
		key       []byte
	}{
		{
			name:      "Should err if key is less than 32",
			key:       xrf197ilz35aq0.RandomBytes(31),
			shouldErr: true,
		},
		{
			name:      "Should decrypt plaintext",
			shouldErr: false,
			key:       encryptionKey,
		},
		{
			name:      "Should err of key is nil",
			shouldErr: true,
			key:       nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decryptedData, err := Decrypt(encryptedData, test.key)
			if test.shouldErr {
				xrf197ilz35aq0.AssertError(t, err)
				assert.Nil(t, decryptedData)
			} else {
				xrf197ilz35aq0.AssertNoError(t, err)
				assert.NotNil(t, decryptedData)
				assert.Equal(t, plaintext, decryptedData)
			}
		})
	}
}

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
