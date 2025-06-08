package config

import (
	"errors"
	"time"
)

// AuthConfig holds configuration for authentication
type AuthConfig struct {
	JWTSecret            string
	JWTExpirationHours   int
	PasswordSalt         string
	VerificationCodeLen  int
	ResetCodeLen         int
	ResetExpirationHours int
}

// loadAuthConfig loads authentication configuration from environment variables
func loadAuthConfig() AuthConfig {
	return AuthConfig{
		JWTSecret:            GetEnv("JWT_SECRET", "your-secret-key"),
		JWTExpirationHours:   GetEnvInt("JWT_EXPIRATION_HOURS", 24),
		PasswordSalt:         GetEnv("PASSWORD_SALT", "your-salt-here"),
		VerificationCodeLen:  GetEnvInt("VERIFICATION_CODE_LEN", 6),
		ResetCodeLen:         GetEnvInt("RESET_CODE_LEN", 6),
		ResetExpirationHours: GetEnvInt("RESET_EXPIRATION_HOURS", 24),
	}
}

// Validate validates the authentication configuration
func (c AuthConfig) Validate() error {
	if c.JWTSecret == "your-secret-key" {
		return errors.New("JWT_SECRET is using the default value, please set a secure secret")
	}
	if c.PasswordSalt == "your-salt-here" {
		return errors.New("PASSWORD_SALT is using the default value, please set a secure salt")
	}
	return nil
}

// JWTExpiration returns the JWT expiration duration
func (c AuthConfig) JWTExpiration() time.Duration {
	return time.Duration(c.JWTExpirationHours) * time.Hour
}

// ResetExpiration returns the password reset expiration duration
func (c AuthConfig) ResetExpiration() time.Duration {
	return time.Duration(c.ResetExpirationHours) * time.Hour
}
