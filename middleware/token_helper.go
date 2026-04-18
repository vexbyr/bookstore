package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a signed JWT for the given userID.
// This is exposed so main.go can print a test token on startup —
// it should NOT be a public API endpoint in production.
func GenerateToken(userID int) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
