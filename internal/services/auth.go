package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dinis/musync/internal/config"
	"github.com/dinis/musync/internal/database"
	"github.com/dinis/musync/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

// AuthService is a simple struct for authentication operations
type AuthService struct {
	db          *database.DB
	config      config.AuthConfig
	emailConfig config.EmailConfig
	emailSvc    *EmailService
}

// NewAuthService creates a new AuthService
func NewAuthService(db *database.DB, authCfg config.AuthConfig, emailCfg config.EmailConfig) *AuthService {
	emailSvc := NewEmailService(emailCfg)
	return &AuthService{
		db:          db,
		config:      authCfg,
		emailConfig: emailCfg,
		emailSvc:    emailSvc,
	}
}

// SignUp registers a new user without a password
func (s *AuthService) SignUp(ctx context.Context, email, username string) error {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where(ctx, "email = ? OR username = ?", email, username).First(ctx, &existingUser); err == nil {
		return ErrUserAlreadyExists
	}

	// Generate verification code
	verificationCode := s.GenerateRandomCode()

	user := models.User{
		Email:            email,
		Username:         username,
		VerificationCode: verificationCode,
		IsPasswordSet:    false,
	}

	if err := s.db.Create(ctx, &user); err != nil {
		return err
	}

	// Send verification email with link
	if err := s.emailSvc.SendVerificationEmail(email, username, verificationCode); err != nil {
		// Log the error but don't fail the signup process
		// In a production environment, you might want to handle this differently
		// For example, you might want to delete the user and return an error
		// or queue the email for retry
		// For now, we'll just log the error
		fmt.Printf("Failed to send verification email: %v\n", err)
	}

	return nil
}

// SetPassword sets a user's password during email verification
func (s *AuthService) SetPassword(ctx context.Context, code, password string) error {
	var user models.User
	if err := s.db.Where(ctx, "verification_code = ?", code).First(ctx, &user); err != nil {
		return ErrInvalidVerificationCode
	}

	// Hash password
	passwordHash := s.HashPassword(password)

	user.PasswordHash = passwordHash
	user.IsPasswordSet = true
	user.IsEmailVerified = true
	user.VerificationCode = ""

	if err := s.db.Save(ctx, &user); err != nil {
		return err
	}

	return nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	var user models.User
	if err := s.db.Where(ctx, "email = ?", email).First(ctx, &user); err != nil {
		return "", ErrInvalidCredentials
	}

	if !user.IsEmailVerified {
		return "", ErrEmailNotVerified
	}

	if !user.IsPasswordSet {
		return "", ErrPasswordNotSet
	}

	if !s.VerifyPassword(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyEmail verifies a user's email using the verification code
func (s *AuthService) VerifyEmail(ctx context.Context, code string) error {
	if code == "" {
		return ErrVerificationCodeRequired
	}

	var user models.User
	if err := s.db.Where(ctx, "verification_code = ?", code).First(ctx, &user); err != nil {
		return ErrInvalidVerificationCode
	}

	user.IsEmailVerified = true
	if err := s.db.Save(ctx, &user); err != nil {
		return err
	}

	return nil
}

// RequestPasswordReset initiates a password reset process
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	var user models.User
	if err := s.db.Where(ctx, "email = ?", email).First(ctx, &user); err != nil {
		// Don't reveal if user exists or not
		return nil
	}

	// Generate reset code
	resetCode := s.GenerateRandomCode()
	expiresAt := time.Now().Add(s.config.ResetExpiration())

	user.ResetCode = resetCode
	user.ResetExpiresAt = &expiresAt
	if err := s.db.Save(ctx, &user); err != nil {
		return err
	}

	// Send reset email
	if err := s.emailSvc.SendPasswordResetEmail(email, resetCode); err != nil {
		// Log the error but don't fail the reset process
		// In a production environment, you might want to handle this differently
		fmt.Printf("Failed to send password reset email: %v\n", err)
	}

	return nil
}

// ResetPassword resets a user's password using the reset code
func (s *AuthService) ResetPassword(ctx context.Context, code, newPassword string) error {
	var user models.User
	if err := s.db.Where(ctx, "reset_code = ? AND reset_expires_at > ?", code, time.Now()).First(ctx, &user); err != nil {
		return ErrInvalidResetCode
	}

	user.PasswordHash = s.HashPassword(newPassword)
	user.ResetCode = ""
	user.ResetExpiresAt = nil
	if err := s.db.Save(ctx, &user); err != nil {
		return err
	}

	return nil
}

// HashPassword hashes a password using Argon2id
func (s *AuthService) HashPassword(password string) string {
	// Using Argon2id for password hashing
	salt := []byte(s.config.PasswordSalt)
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(hash)
}

// VerifyPassword verifies a password against a hash
func (s *AuthService) VerifyPassword(password, hash string) bool {
	salt := []byte(s.config.PasswordSalt)
	hashBytes, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false
	}

	passwordHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return string(passwordHash) == string(hashBytes)
}

// GenerateToken generates a JWT token for a user
func (s *AuthService) GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.config.JWTExpiration()).Unix(),
	})

	return token.SignedString([]byte(s.config.JWTSecret))
}

// GenerateRandomCode generates a random code for email verification or password reset
func (s *AuthService) GenerateRandomCode() string {
	codeLength := s.config.VerificationCodeLen
	if codeLength <= 0 {
		codeLength = 6
	}

	// Generate a random code using crypto/rand
	b := make([]byte, codeLength)
	if _, err := rand.Read(b); err != nil {
		// Fallback to a simple code if random generation fails
		return "123456"
	}

	// Convert to a base64 string and take the first codeLength characters
	return base64.StdEncoding.EncodeToString(b)[:codeLength]
}
