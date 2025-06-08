package services

import (
	"context"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/dinis/musync/internal/database"
	"github.com/dinis/musync/internal/models"
)

// FileStorageService handles operations related to file storage and retrieval
type FileStorageService struct {
	db *database.DB
}

// NewFileStorageService creates a new FileStorageService
func NewFileStorageService(db *database.DB) *FileStorageService {
	return &FileStorageService{
		db: db,
	}
}

// GetFileStream returns a reader for a track's file
func (s *FileStorageService) GetFileStream(ctx context.Context, userID, trackID uint) (io.ReadCloser, string, error) {
	// Get the track
	var track models.Track
	if err := s.db.First(ctx, &track, trackID); err != nil {
		return nil, "", errors.New("track not found")
	}

	// Check if the track belongs to a library owned by the user
	var library models.MusicLibrary
	if err := s.db.First(ctx, &library, track.LibraryID); err != nil {
		return nil, "", errors.New("library not found")
	}

	if library.UserID != userID {
		return nil, "", errors.New("unauthorized access to track")
	}

	// Parse the location to determine the source
	location := track.Location
	contentType := s.getContentTypeFromLocation(location)

	// Handle different storage sources
	if strings.HasPrefix(location, "file://localhost") {
		// Local file
		return s.getLocalFileStream(location)
	} else if strings.HasPrefix(location, "pcloud://") {
		// pCloud file
		return s.getPCloudFileStream(location)
	} else {
		// Unsupported source
		return nil, contentType, errors.New("unsupported file location")
	}
}

// getLocalFileStream returns a reader for a local file
func (s *FileStorageService) getLocalFileStream(location string) (io.ReadCloser, string, error) {
	// Remove the file:// prefix and decode URL
	path := strings.TrimPrefix(location, "file://localhost")
	path, err := url.PathUnescape(path)
	if err != nil {
		return nil, "", errors.New("invalid file path")
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, "", errors.New("failed to open file: " + err.Error())
	}

	contentType := s.getContentTypeFromLocation(location)
	return file, contentType, nil
}

// getPCloudFileStream returns a reader for a pCloud file
func (s *FileStorageService) getPCloudFileStream(location string) (io.ReadCloser, string, error) {
	// In a real implementation, this would use the pCloud API to get a download link
	// and then return a reader for that link
	// For now, we'll just return an error
	return nil, "", errors.New("pCloud integration not implemented yet")

	// Example of how this might be implemented:
	/*
		// Extract the file ID from the pCloud URL
		fileID := strings.TrimPrefix(location, "pcloud://")

		// Use the pCloud API to get a download link
		downloadLink, err := s.pCloudClient.GetFileLink(fileID)
		if err != nil {
			return nil, "", errors.New("failed to get pCloud download link: " + err.Error())
		}

		// Make an HTTP request to get the file
		resp, err := http.Get(downloadLink)
		if err != nil {
			return nil, "", errors.New("failed to download file from pCloud: " + err.Error())
		}

		contentType := resp.Header.Get("Content-Type")
		return resp.Body, contentType, nil
	*/
}

// getContentTypeFromLocation returns the content type based on the file extension
func (s *FileStorageService) getContentTypeFromLocation(location string) string {
	ext := strings.ToLower(filepath.Ext(location))
	switch ext {
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".flac":
		return "audio/flac"
	case ".aac":
		return "audio/aac"
	case ".ogg":
		return "audio/ogg"
	default:
		return "application/octet-stream"
	}
}
