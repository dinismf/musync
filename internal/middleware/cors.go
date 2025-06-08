package middleware

import (
	"github.com/dinis/musync/internal/config"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware is a middleware for CORS
type CorsMiddleware struct {
	config config.ServerConfig
}

// NewCorsMiddleware creates a new CorsMiddleware
func NewCorsMiddleware(cfg config.ServerConfig) *CorsMiddleware {
	return &CorsMiddleware{
		config: cfg,
	}
}

// Cors is a middleware that handles CORS
func (m *CorsMiddleware) Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Origin", m.config.AllowedOrigins[0])
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}