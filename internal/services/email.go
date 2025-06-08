package services

import (
	"fmt"
	"net/smtp"

	"github.com/dinis/musync/internal/config"
)

// EmailService handles sending emails
type EmailService struct {
	config config.EmailConfig
}

// NewEmailService creates a new EmailService
func NewEmailService(cfg config.EmailConfig) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

// SendVerificationEmail sends an email with a verification link
func (s *EmailService) SendVerificationEmail(to, username, code string) error {
	subject := "Verify Your Musync Account"
	verificationURL := fmt.Sprintf("%s/verify-email?code=%s", s.config.FrontendURL, code)

	body := fmt.Sprintf(`
	<html>
	<body>
		<h2>Welcome to Musync, %s!</h2>
		<p>Thank you for registering. Please click the link below to verify your email and set your password:</p>
		<p><a href="%s">Verify Email and Set Password</a></p>
		<p>If you did not register for a Musync account, please ignore this email.</p>
	</body>
	</html>
	`, username, verificationURL)

	return s.sendEmail(to, subject, body)
}

// SendPasswordResetEmail sends an email with a password reset link
func (s *EmailService) SendPasswordResetEmail(to, code string) error {
	subject := "Reset Your Musync Password"
	resetURL := fmt.Sprintf("%s/reset-password?code=%s", s.config.FrontendURL, code)

	body := fmt.Sprintf(`
	<html>
	<body>
		<h2>Password Reset Request</h2>
		<p>You requested to reset your password. Please click the link below to set a new password:</p>
		<p><a href="%s">Reset Password</a></p>
		<p>If you did not request a password reset, please ignore this email.</p>
	</body>
	</html>
	`, resetURL)

	return s.sendEmail(to, subject, body)
}

// sendEmail sends an email using SMTP
func (s *EmailService) sendEmail(to, subject, body string) error {
	// Construct email headers
	headers := make(map[string]string)
	headers["From"] = s.config.From
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Construct message
	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	// Prepare SMTP connection details
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// For Mailpit in local development, we don't need authentication
	// But for production, we would use authentication
	var auth smtp.Auth
	if s.config.Username != "" && s.config.Password != "" {
		auth = smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	}

	// Send email
	err := smtp.SendMail(addr, auth, s.config.From, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
