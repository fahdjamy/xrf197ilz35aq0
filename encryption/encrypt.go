package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// Generate a 32-byte (256-bit) key, suitable for AES-256
const minKeySize = 32

// Larger keys can lead to slightly slower encryption and decryption operations.
const maxKeySize = 256

// https://docs.google.com/document/d/1uqD8gAjpAN4EWsmg7yv1AbcKL_8lJxXGNTfO70YII_0

// Encrypt has Separate Nonces: We generate gcmNonce and aadNonce separately.
func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	// Create the AES cipher block using the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a GCM cipher using the block
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate a nonce for the GCM encryption
	gcmNonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, gcmNonce); err != nil {
		return nil, err
	}

	// Generate a separate nonce for the additional authenticated data (AAD) Additional Authenticated Data
	aadNonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, aadNonce); err != nil {
		return nil, err
	}

	// Encrypt the plaintext using the GCM, providing both nonces and any additional data
	ciphertext := aesgcm.Seal(nil, gcmNonce, plaintext, aadNonce)
	return append(gcmNonce, append(aadNonce, ciphertext...)...), nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < 2*nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	gcmNonce := ciphertext[:nonceSize]
	aadNonce := ciphertext[nonceSize : 2*nonceSize]
	ciphertext = ciphertext[2*nonceSize:]

	plaintext, err := aesgcm.Open(nil, gcmNonce, ciphertext, aadNonce)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func GenerateKey(keySize int) ([]byte, error) {
	// keySize. For strong encryption, we recommend using at least 32 bytes (256 bits) for AES-256.
	// Create a byte slice to hold the key
	if keySize < minKeySize || keySize > maxKeySize {
		return nil, fmt.Errorf("key size generated should be between [%d, %d]", minKeySize, maxKeySize)
	}
	key := make([]byte, keySize)

	// reads random bytes from the cryptographically secure random source (crypto/rand) and
	// fills the key slice with these random bytes.
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	return key, nil
}
