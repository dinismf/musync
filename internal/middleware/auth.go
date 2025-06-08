package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dinis/musync/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware is a middleware for authentication
type AuthMiddleware struct {
	config config.AuthConfig
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(cfg config.AuthConfig) *AuthMiddleware {
	return &AuthMiddleware{
		config: cfg,
	}
}

// RequireAuth is a middleware that requires authentication
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Parse the token
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(m.config.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if the token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Add the user ID to the context
			userID, ok := claims["user_id"].(float64)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				c.Abort()
				return
			}
			c.Set("user_id", uint(userID))
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
}