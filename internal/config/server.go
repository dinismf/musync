package config

// ServerConfig holds configuration for the server
type ServerConfig struct {
	Port            string
	AllowedOrigins  []string
	AllowedMethods  []string
	AllowedHeaders  []string
	TrustedProxies  []string
	RateLimitEnable bool
	RateLimitMax    int
}

// loadServerConfig loads server configuration from environment variables
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:           GetEnv("PORT", "8080"),
		AllowedOrigins: []string{GetEnv("CORS_ALLOWED_ORIGINS", "*")},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
		},
		AllowedHeaders: []string{
			"Content-Type", "Authorization", "X-Requested-With",
		},
		TrustedProxies:  []string{GetEnv("TRUSTED_PROXIES", "127.0.0.1")},
		RateLimitEnable: GetEnvBool("RATE_LIMIT_ENABLE", true),
		RateLimitMax:    GetEnvInt("RATE_LIMIT_MAX", 100),
	}
}
