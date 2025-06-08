package config

// EmailConfig holds configuration for email sending
type EmailConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	From        string
	FrontendURL string
}

// loadEmailConfig loads email configuration from environment variables
func loadEmailConfig() EmailConfig {
	return EmailConfig{
		Host:        GetEnv("EMAIL_HOST", "localhost"),
		Port:        GetEnvInt("EMAIL_PORT", 1025),
		Username:    GetEnv("EMAIL_USERNAME", ""),
		Password:    GetEnv("EMAIL_PASSWORD", ""),
		From:        GetEnv("EMAIL_FROM", "noreply@musync.com"),
		FrontendURL: GetEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}

// Validate validates the email configuration
func (c *EmailConfig) Validate() error {
	// For local development with Mailpit, we don't need to validate credentials
	return nil
}
