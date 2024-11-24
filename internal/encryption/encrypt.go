package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	xrfErr "xrf197ilz35aq0/internal/error"
)

// Generate a 32-byte (256-bit) key, suitable for AES-256
const minKeySize = 32

// Larger keys can lead to slightly slower encryption and decryption operations.
const maxKeySize = 256

// https://docs.google.com/document/d/1uqD8gAjpAN4EWsmg7yv1AbcKL_8lJxXGNTfO70YII_0

// Encrypt has Separate Nonces: We generate gcmNonce and aadNonce separately.
func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	if len(plaintext) < 3 {
		return nil, &Error{
			message: "plaintext should at least be of length 3",
		}
	}
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, &Error{
			message: "Invalid key size. Key must be 16, 24, or 32 bytes",
		}
	}
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
	if len(key) < minKeySize {
		return nil, &Error{
			message: "for security, encryption key should at least be 32 bytes",
		}
	}

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
		if err == io.EOF {
			// io.EOF: This error is returned if the underlying source of randomness (e.g., /dev/urandom)
			// unexpectedly reaches its end. While rare, it's possible in scenarios where the system is under
			// extreme stress or there's an issue with the entropy source.
			fmt.Println("Unexpected end of randomness source")
			return nil, &xrfErr.Internal{
				Message: "Unexpected end of randomness source",
				Source:  "GenerateKey",
				Err:     err,
			}
		}
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	return key, nil
}
