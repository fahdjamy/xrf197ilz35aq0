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

func TestEncryptAndDecodeInvalidInputs(t *testing.T) {
	key, err := encryption.GenerateKey(32)
	internal.AssertNoError(t, err)
	plaintext := internal.RandomBytes(5)

	// Test with an invalid key size
	_, err = encryption.EncryptAndEncode(plaintext, []byte("invalid key"))
	internal.AssertError(t, err)

	// Test with invalid encoded ciphertext (e.g., not Base64)
	_, err = encryption.DecodeAndDecrypt("invalid encoded ciphertext", key)
	internal.AssertError(t, err)

	// Test with empty plaintext
	_, err = encryption.EncryptAndEncode([]byte(""), key)
	internal.AssertError(t, err)

	// Test with valid values
	encryptedData, err := encryption.EncryptAndEncode(plaintext, key)
	internal.AssertNoError(t, err)
	_, err = encryption.DecodeAndDecrypt(encryptedData, key)
	internal.AssertNoError(t, err)
}
