package services

import (
	"errors"
)

// Service error definitions
var (
	// Auth service errors
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrEmailNotVerified         = errors.New("email not verified")
	ErrPasswordNotSet           = errors.New("password not set")
	ErrVerificationCodeRequired = errors.New("verification code is required")
	ErrInvalidVerificationCode  = errors.New("invalid verification code")
	ErrInvalidResetCode         = errors.New("invalid or expired reset code")
)
