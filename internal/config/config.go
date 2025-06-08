package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Redis    RedisConfig
	Email    EmailConfig
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		Server:   loadServerConfig(),
		Database: loadDatabaseConfig(),
		Auth:     loadAuthConfig(),
		Redis:    loadRedisConfig(),
		Email:    loadEmailConfig(),
	}

	return cfg, nil
}

// GetEnv gets an environment variable or returns a default value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetEnvInt gets an environment variable as an integer or returns a default value
func GetEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid value for %s, using default: %v\n", key, defaultValue)
		return defaultValue
	}

	return value
}

// GetEnvBool gets an environment variable as a boolean or returns a default value
func GetEnvBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid value for %s, using default: %v\n", key, defaultValue)
		return defaultValue
	}

	return value
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config validation failed: %w", err)
	}

	if err := c.Auth.Validate(); err != nil {
		return fmt.Errorf("auth config validation failed: %w", err)
	}

	if err := c.Email.Validate(); err != nil {
		return fmt.Errorf("email config validation failed: %w", err)
	}

	return nil
}
