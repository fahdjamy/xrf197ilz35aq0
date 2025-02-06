package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"xrf197ilz35aq0/internal/encryption"
)

type UserTokenPayload struct {
	UserId    string
	ServerId  string
	ExpiresAt time.Duration
}

func GenerateAuthToken(payload UserTokenPayload, secret []byte) (string, error) {
	// High Entropy is Key (for the input to the HMAC, not the hash itself):  Whether you use SHA-256 directly (not recommended)
	//or HMAC-SHA256 (recommended), the security of the token fundamentally depends on the input having high entropy.
	//This means the input (the secret key plus any other data you include) must be sufficiently random and long
	//to prevent attackers from guessing it.  Think of it this way: SHA-256 is like a very strong blender;
	//it makes anything you put in hard to recognize. But if you put in a tiny piece of paper with a short,
	//predictable word, someone could still figure out what it was by trying every short word.

	// ::: serialize the token payload
	now := time.Now()
	payloadBytes, err := json.Marshal(struct {
		UserId    string
		ServerId  string
		ExpiresAt time.Time
		IssuedAt  time.Time
	}{
		IssuedAt:  now,
		UserId:    payload.UserId,
		ServerId:  payload.ServerId,
		ExpiresAt: now.Add(payload.ExpiresAt),
	})

	if err != nil {
		return "", fmt.Errorf("error marshaling payload: %w", err)
	}

	// ::: generate salt. generating a salt is similar to generating a secret key except, we are keeping the size small
	salt, err := encryption.GenerateKey(21)
	if err != nil {
		return "", fmt.Errorf("error generating salt: %w", err)
	}

	// ::: Combine token Payload, and salt.
	data := append(payloadBytes, salt...)

	// ::: Create an HMAC-SHA256 signature using SHA-256 to sign the data. pass in a secret key
	h := hmac.New(sha256.New, secret)
	_, err = h.Write(data)
	if err != nil {
		return "", fmt.Errorf("error hashing payload: %w", err)
	}

	// Avoid the standard Base64 encoding unless you're certain the tokens will never be used in a URL context, and you
	//	manually handle the encoding of the + and / characters if they appear base64.URLEncoding.EncodeToString is
	//	generally the best choice for authentication tokens due to its balance of compactness and URL safety
	token := base64.StdEncoding.EncodeToString(h.Sum(nil))
	// If you absolutely want to avoid any potential URL encoding issues and the larger token size is acceptable,
	// hexadecimal encoding is a safe bet: is a good alternative if URL safety is paramount & the larger token size isn't an issue
	// token := hex.EncodeToString(h.Sum(nil))

	return token, nil
}

func generateHash() {

}
