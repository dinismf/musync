package main

import (
	"github.com/dinis/musync/internal/config"
	"github.com/dinis/musync/internal/database"
	"github.com/dinis/musync/internal/logging"
	"github.com/dinis/musync/internal/middleware"
	"github.com/dinis/musync/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logging
	logging.Init(logging.InfoLevel, nil)
	logger := logging.GetLogger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logger.Fatal("Configuration validation failed: %v", err)
	}

	// Initialize database
	database.InitDB(cfg.Database)

	// Set up Gin router
	r := gin.New() // Use gin.New() instead of gin.Default() to avoid using the default logger and recovery middleware

	// Setup middleware
	errorMiddleware := middleware.NewErrorMiddleware()
	r.Use(errorMiddleware.ErrorHandler()) // Add error handling middleware
	r.Use(gin.Recovery()) // Add recovery middleware to handle panics

	// Setup logging middleware
	loggingMiddleware := middleware.NewLoggingMiddleware()
	r.Use(loggingMiddleware.Logger()) // Add custom logger middleware

	// Setup CORS
	corsMiddleware := middleware.NewCorsMiddleware(cfg.Server)
	r.Use(corsMiddleware.Cors())

	// Setup routes
	routes.SetupRoutes(r, cfg)

	// Start the server
	logger.Info("Starting server on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Failed to start server: %v", err)
	}
}
