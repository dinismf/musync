package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dinis/musync/internal/database"
	"github.com/dinis/musync/internal/dto"
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

	// Convert to response DTOs
	libraryResponses := dto.ToLibraryResponses(libraries)
	c.JSON(http.StatusOK, libraryResponses)
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

	// Convert to response DTO
	libraryResponse := dto.ToLibraryResponse(*library)
	c.JSON(http.StatusOK, libraryResponse)
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

	// Convert to response DTOs
	trackResponses := dto.ToTrackResponses(tracks)
	c.JSON(http.StatusOK, trackResponses)
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

	// Convert to response DTOs
	playlistResponses := dto.ToPlaylistResponses(playlists)
	c.JSON(http.StatusOK, playlistResponses)
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

	// Convert to response DTOs
	trackResponses := dto.ToTrackResponses(tracks)
	c.JSON(http.StatusOK, trackResponses)
}

// StreamTrack streams a track's audio file or returns track information for local files
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

	// Stream the file through the backend
	// This ensures compatibility with web browsers that can't access local files directly
	fileStream, contentType, err := h.fileService.GetFileStream(c.Request.Context(), userID.(uint), uint(trackID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get track: " + err.Error()})
		return
	}
	defer fileStream.Close()

	// Set content type header
	c.Header("Content-Type", contentType)

	// Add headers to enable seeking and other audio controls
	c.Header("Accept-Ranges", "bytes")
	c.Header("Cache-Control", "no-cache")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Check if this is a range request
	rangeHeader := c.Request.Header.Get("Range")
	if rangeHeader != "" {
		// Parse the range header
		rangeStr := strings.Replace(rangeHeader, "bytes=", "", 1)
		rangeParts := strings.Split(rangeStr, "-")

		// Get the start position
		startPos, err := strconv.ParseInt(rangeParts[0], 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid range header"})
			return
		}

		// Get file size
		fileInfo, err := fileStream.(interface{ Stat() (os.FileInfo, error) }).Stat()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
			return
		}
		fileSize := fileInfo.Size()

		// Calculate end position
		endPos := fileSize - 1
		if len(rangeParts) > 1 && rangeParts[1] != "" {
			endPos, err = strconv.ParseInt(rangeParts[1], 10, 64)
			if err != nil {
				endPos = fileSize - 1
			}
		}

		// Ensure endPos is not greater than fileSize
		if endPos >= fileSize {
			endPos = fileSize - 1
		}

		// Calculate content length
		contentLength := endPos - startPos + 1

		// Seek to the requested position
		_, err = fileStream.Seek(startPos, io.SeekStart)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to seek to position"})
			return
		}

		// Set headers for partial content
		c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", startPos, endPos, fileSize))
		c.Header("Content-Length", fmt.Sprintf("%d", contentLength))
		c.Status(http.StatusPartialContent)
	} else {
		// Full content request
		c.Status(http.StatusOK)
	}

	// Stream the file
	_, err = io.Copy(c.Writer, fileStream)
	if err != nil {
		// The response has already started, so we can't return a JSON error
		// Just log the error and return
		c.Status(http.StatusInternalServerError)
		return
	}
}

// DeleteLibrary handles the deletion of a music library and all its associated resources
func (h *MusicLibraryHandler) DeleteLibrary(c *gin.Context) {
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

	// Delete the library
	err = h.libraryService.DeleteLibrary(c.Request.Context(), userID.(uint), uint(libraryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete library: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Library deleted successfully"})
}
