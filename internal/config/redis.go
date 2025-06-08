package config

import "fmt"

// RedisConfig holds configuration for Redis
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Enabled  bool
}

// loadRedisConfig loads Redis configuration from environment variables
func loadRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     GetEnv("REDIS_HOST", "localhost"),
		Port:     GetEnv("REDIS_PORT", "6379"),
		Password: GetEnv("REDIS_PASSWORD", ""),
		DB:       GetEnvInt("REDIS_DB", 0),
		Enabled:  GetEnvBool("REDIS_ENABLED", false),
	}
}

// DSN returns the Redis connection string
func (c RedisConfig) DSN() string {
	if c.Password == "" {
		return fmt.Sprintf("%s:%s", c.Host, c.Port)
	}
	return fmt.Sprintf("redis://:%s@%s:%s", c.Password, c.Host, c.Port)
}
