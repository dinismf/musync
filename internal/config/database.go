package config

import (
	"errors"
	"fmt"
)

// DatabaseConfig holds configuration for the database
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	MaxOpen  int
	MaxIdle  int
	Lifetime int // in seconds
}

// loadDatabaseConfig loads database configuration from environment variables
func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     GetEnv("DB_HOST", "localhost"),
		Port:     GetEnv("DB_PORT", "5432"),
		User:     GetEnv("DB_USER", "postgres"),
		Password: GetEnv("DB_PASSWORD", "postgres"),
		Name:     GetEnv("DB_NAME", "musync"),
		SSLMode:  GetEnv("DB_SSLMODE", "disable"),
		MaxOpen:  GetEnvInt("DB_MAX_OPEN", 25),
		MaxIdle:  GetEnvInt("DB_MAX_IDLE", 5),
		Lifetime: GetEnvInt("DB_LIFETIME", 300),
	}
}

// DSN returns the database connection string
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// Validate validates the database configuration
func (c DatabaseConfig) Validate() error {
	if c.Host == "" {
		return errors.New("database host is required")
	}
	if c.Port == "" {
		return errors.New("database port is required")
	}
	if c.User == "" {
		return errors.New("database user is required")
	}
	if c.Name == "" {
		return errors.New("database name is required")
	}
	return nil
}
