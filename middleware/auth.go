package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ContextKeyUserID is the key used to store the authenticated user ID in the
// Gin context so handlers can retrieve it with c.MustGet(ContextKeyUserID).
const ContextKeyUserID = "userID"

// jwtSecret is the HMAC secret used to sign and verify tokens.
// In production, load this from an environment variable or secrets manager.
var jwtSecret = []byte("super-secret-key-change-in-production")

// Claims defines the JWT payload we expect.
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// RequireAuth is a Gin middleware that:
//  1. Reads the Authorization: Bearer <token> header.
//  2. Parses and validates the JWT (signature + expiry).
//  3. Stores the user_id claim in the Gin context for downstream handlers.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		// Expect "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be: Bearer <token>"})
			return
		}

		tokenStr := parts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			// Enforce HMAC signing method to prevent the "alg: none" attack.
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		if claims.UserID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token missing user_id claim"})
			return
		}

		// Inject user_id into the context for handlers.
		c.Set(ContextKeyUserID, claims.UserID)
		c.Next()
	}
}

// GetUserID is a helper that extracts the authenticated user ID from the
// Gin context. Panics if called outside of a RequireAuth-protected route.
func GetUserID(c *gin.Context) int {
	return c.MustGet(ContextKeyUserID).(int)
}
