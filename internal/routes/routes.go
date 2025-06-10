package routes

import (
	"github.com/dinis/musync/internal/config"
	"github.com/dinis/musync/internal/handlers"
	"github.com/dinis/musync/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	// Add a root handler
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Musync API",
		})
	})

	// Add a health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	authHandler := handlers.NewAuthHandler(cfg.Auth, cfg.Email)
	musicLibraryHandler := handlers.NewMusicLibraryHandler()

	// Public routes
	public := r.Group("/api")
	{
		// Auth routes
		auth := public.Group("/auth")
		{
			auth.POST("/signup", authHandler.SignUp)
			auth.POST("/login", authHandler.Login)
			auth.GET("/verify", authHandler.VerifyEmail)
			auth.POST("/set-password", authHandler.SetPassword)
			auth.POST("/reset-password", authHandler.RequestPasswordReset)
			auth.POST("/confirm-reset", authHandler.ResetPassword)
		}
	}

	// Protected routes
	protected := r.Group("/api")
	{
		// Add middleware for authentication
		authMiddleware := middleware.NewAuthMiddleware(cfg.Auth)
		protected.Use(authMiddleware.RequireAuth())

		// Music library routes
		library := protected.Group("/libraries")
		{
			library.POST("", musicLibraryHandler.UploadLibrary)
			library.GET("", musicLibraryHandler.GetLibraries)
			library.GET("/:id", musicLibraryHandler.GetLibrary)
			library.GET("/:id/tracks", musicLibraryHandler.GetTracks)
			library.GET("/:id/playlists", musicLibraryHandler.GetPlaylists)
			library.DELETE("/:id", musicLibraryHandler.DeleteLibrary)
		}

		// Playlist routes
		playlist := protected.Group("/playlists")
		{
			playlist.GET("/:id/tracks", musicLibraryHandler.GetPlaylistTracks)
		}

		// Track routes
		track := protected.Group("/tracks")
		{
			track.GET("/:id/stream", musicLibraryHandler.StreamTrack)
		}

		// The following route groups are commented out to avoid unused variable warnings
		// They are left here as a template for future implementation

		// Example of how to set up protected routes:
		// protected.GET("/me", userHandler.GetMe)

		/*
			// User routes
			user := protected.Group("/user")
			{
				user.GET("/profile", userHandler.GetProfile)
				user.PUT("/profile", userHandler.UpdateProfile)
			}

			// Artist routes
			artist := protected.Group("/artist")
			{
				artist.POST("", artistHandler.Create)
				artist.GET("/:id", artistHandler.Get)
				artist.PUT("/:id", artistHandler.Update)
				artist.DELETE("/:id", artistHandler.Delete)
			}

			// Label routes
			label := protected.Group("/label")
			{
				label.POST("", labelHandler.Create)
				label.GET("/:id", labelHandler.Get)
				label.PUT("/:id", labelHandler.Update)
				label.DELETE("/:id", labelHandler.Delete)
			}

			// Release routes
			release := protected.Group("/release")
			{
				release.POST("", releaseHandler.Create)
				release.GET("/:id", releaseHandler.Get)
				release.PUT("/:id", releaseHandler.Update)
				release.DELETE("/:id", releaseHandler.Delete)
			}
		*/
	}
}
