package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Standard error types
var (
	ErrNotFound          = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrBadRequest        = errors.New("bad request")
	ErrInternalServer    = errors.New("internal server error")
	ErrValidation        = errors.New("validation error")
	ErrDuplicateResource = errors.New("resource already exists")
	ErrTimeout           = errors.New("operation timed out")
)

// AppError represents an application error
type AppError struct {
	Err        error
	StatusCode int
	Message    string
	Details    interface{}
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(err error, statusCode int, message string, details interface{}) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, details interface{}) *AppError {
	return NewAppError(ErrNotFound, http.StatusNotFound, message, details)
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string, details interface{}) *AppError {
	return NewAppError(ErrUnauthorized, http.StatusUnauthorized, message, details)
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string, details interface{}) *AppError {
	return NewAppError(ErrForbidden, http.StatusForbidden, message, details)
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string, details interface{}) *AppError {
	return NewAppError(ErrBadRequest, http.StatusBadRequest, message, details)
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(message string, details interface{}) *AppError {
	return NewAppError(ErrInternalServer, http.StatusInternalServerError, message, details)
}

// NewValidationError creates a new validation error
func NewValidationError(message string, details interface{}) *AppError {
	return NewAppError(ErrValidation, http.StatusBadRequest, message, details)
}

// NewDuplicateResourceError creates a new duplicate resource error
func NewDuplicateResourceError(message string, details interface{}) *AppError {
	return NewAppError(ErrDuplicateResource, http.StatusConflict, message, details)
}

// NewTimeoutError creates a new timeout error
func NewTimeoutError(message string, details interface{}) *AppError {
	return NewAppError(ErrTimeout, http.StatusRequestTimeout, message, details)
}

// Is checks if the target error is of the same type as the source error
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Wrap wraps an error with a message
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}