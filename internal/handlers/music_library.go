package handlers

import (
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/dinis/musync/internal/database"
	"github.com/dinis/musync/internal/services"
	"github.com/gin-gonic/gin"
)

// MusicLibraryHandler handles HTTP requests related to music libraries
type MusicLibraryHandler struct {
	libraryService *services.MusicLibraryService
	fileService    *services.FileStorageService
}

// NewMusicLibraryHandler creates a new MusicLibraryHandler
func NewMusicLibraryHandler() *MusicLibraryHandler {
	libraryService := services.NewMusicLibraryService(database.GlobalDB)
	fileService := services.NewFileStorageService(database.GlobalDB)
	return &MusicLibraryHandler{
		libraryService: libraryService,
		fileService:    fileService,
	}
}

// UploadLibraryRequest represents the JSON request for uploading a library
type UploadLibraryRequest struct {
	Name     string `json:"name" binding:"required"`
	FileData string `json:"file_data" binding:"required"` // Base64 encoded XML file
}

// UploadLibrary handles the upload of a music library XML file
func (h *MusicLibraryHandler) UploadLibrary(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse JSON request
	var req UploadLibraryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate library name
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Library name is required"})
		return
	}

	// Decode base64 file data
	fileData, err := base64.StdEncoding.DecodeString(req.FileData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file data encoding"})
		return
	}

	// Check if the data appears to be XML
	fileDataStr := string(fileData)
	if !strings.HasPrefix(fileDataStr, "<?xml") && !strings.HasPrefix(fileDataStr, "<DJ_PLAYLISTS") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an XML file"})
		return
	}

	// Create a reader from the decoded data
	xmlReader := strings.NewReader(fileDataStr)

	// Upload the library
	libraryID, err := h.libraryService.UploadLibrary(c.Request.Context(), userID.(uint), req.Name, xmlReader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload library: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Library uploaded successfully", "library_id": libraryID})
}

// GetLibraries returns all music libraries for the authenticated user
func (h *MusicLibraryHandler) GetLibraries(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get libraries
	libraries, err := h.libraryService.GetLibraries(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get libraries"})
		return
	}

	c.JSON(http.StatusOK, libraries)
}

// GetLibrary returns a specific music library
func (h *MusicLibraryHandler) GetLibrary(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get library ID from URL
	libraryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid library ID"})
		return
	}

	// Get library
	library, err := h.libraryService.GetLibrary(c.Request.Context(), userID.(uint), uint(libraryID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Library not found"})
		return
	}

	c.JSON(http.StatusOK, library)
}

// GetTracks returns all tracks in a library
func (h *MusicLibraryHandler) GetTracks(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get library ID from URL
	libraryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid library ID"})
		return
	}

	// Get tracks
	tracks, err := h.libraryService.GetTracks(c.Request.Context(), userID.(uint), uint(libraryID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Library not found"})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

// GetPlaylists returns all playlists in a library
func (h *MusicLibraryHandler) GetPlaylists(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get library ID from URL
	libraryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid library ID"})
		return
	}

	// Get playlists
	playlists, err := h.libraryService.GetPlaylists(c.Request.Context(), userID.(uint), uint(libraryID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Library not found"})
		return
	}

	c.JSON(http.StatusOK, playlists)
}

// GetPlaylistTracks returns all tracks in a playlist
func (h *MusicLibraryHandler) GetPlaylistTracks(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get playlist ID from URL
	playlistID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	// Get tracks
	tracks, err := h.libraryService.GetPlaylistTracks(c.Request.Context(), userID.(uint), uint(playlistID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

// StreamTrack streams a track's audio file
func (h *MusicLibraryHandler) StreamTrack(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get track ID from URL
	trackID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	// Get file stream
	fileStream, contentType, err := h.fileService.GetFileStream(c.Request.Context(), userID.(uint), uint(trackID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get track: " + err.Error()})
		return
	}
	defer fileStream.Close()

	// Set content type header
	c.Header("Content-Type", contentType)

	// Stream the file
	c.Status(http.StatusOK)
	_, err = io.Copy(c.Writer, fileStream)
	if err != nil {
		// The response has already started, so we can't return a JSON error
		// Just log the error and return
		c.Status(http.StatusInternalServerError)
		return
	}
}
