package middleware

import (
	"net/http"

	"github.com/dinis/musync/internal/errors"
	"github.com/dinis/musync/internal/logging"
	"github.com/gin-gonic/gin"
)

// ErrorMiddleware is a middleware for error handling
type ErrorMiddleware struct {
	logger *logging.Logger
}

// NewErrorMiddleware creates a new ErrorMiddleware
func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{
		logger: logging.GetLogger(),
	}
}

// ErrorHandler is a middleware that handles errors
func (m *ErrorMiddleware) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last().Err

			// Check if it's an AppError
			var appErr *errors.AppError
			if errors.As(err, &appErr) {
				// Log the error
				m.logger.Error("AppError: %v", appErr)

				// Return the error response
				c.JSON(appErr.StatusCode, gin.H{
					"error":   appErr.Error(),
					"details": appErr.Details,
				})
				return
			}

			// Handle other errors
			m.logger.Error("Error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
	}
}

// Deprecated: Use NewErrorMiddleware().ErrorHandler() instead
// ErrorHandler is a middleware that handles errors
func ErrorHandler() gin.HandlerFunc {
	return NewErrorMiddleware().ErrorHandler()
}
