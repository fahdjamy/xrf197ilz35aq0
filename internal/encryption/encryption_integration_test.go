package encryption_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/encryption"
)

func TestEncryptAndDecrypt(t *testing.T) {
	t.Run("generates key, encrypts and decrypts data", func(t *testing.T) {
		key, err := encryption.GenerateKey(32)
		internal.AssertNoError(t, err)

		plainText := internal.RandomBytes(5)
		encryptedData, err := encryption.Encrypt(plainText, key)
		internal.AssertNoError(t, err)
		assert.NotNil(t, encryptedData)

		decryptedData, err := encryption.Decrypt(encryptedData, key)
		internal.AssertNoError(t, err)
		assert.Equal(t, plainText, decryptedData)
	})
}
