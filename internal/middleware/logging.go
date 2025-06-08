package middleware

import (
	"time"

	"github.com/dinis/musync/internal/logging"
	"github.com/gin-gonic/gin"
)

// LoggingMiddleware is a middleware for logging requests
type LoggingMiddleware struct {
	logger *logging.Logger
}

// NewLoggingMiddleware creates a new LoggingMiddleware
func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logging.GetLogger(),
	}
}

// Logger is a middleware that logs the request
func (m *LoggingMiddleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request details
		m.logger.LogRequest(
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			latency,
		)
	}
}

// Deprecated: Use NewLoggingMiddleware().Logger() instead
// Logger is a middleware that logs the request
func Logger() gin.HandlerFunc {
	return NewLoggingMiddleware().Logger()
}
