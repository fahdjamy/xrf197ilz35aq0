package security

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWTToken(userId string, secret []byte, hours time.Duration, serverId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId":    userId,
			"issuer":    serverId,
			"expiresIn": time.Now().Add(time.Hour * hours).Unix(), // 1-hour expiration
		})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
