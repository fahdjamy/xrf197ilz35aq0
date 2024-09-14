package encryption_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/internal/encryption"
)

func TestEncryptAndDecrypt(t *testing.T) {
	t.Run("generates key, encrypts and decrypts data", func(t *testing.T) {
		key, err := encryption.GenerateKey(32)
		xrf197ilz35aq0.AssertNoError(t, err)

		plainText := xrf197ilz35aq0.RandomBytes(5)
		encryptedData, err := encryption.Encrypt(plainText, key)
		xrf197ilz35aq0.AssertNoError(t, err)
		assert.NotNil(t, encryptedData)

		decryptedData, err := encryption.Decrypt(encryptedData, key)
		xrf197ilz35aq0.AssertNoError(t, err)
		assert.Equal(t, plainText, decryptedData)
	})
}
