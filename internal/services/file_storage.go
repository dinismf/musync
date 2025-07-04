package services

import (
	"context"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dinis/musync/internal/database"
	"github.com/dinis/musync/internal/models"
)

// ReadSeekCloser combines io.ReadSeeker and io.Closer interfaces
type ReadSeekCloser interface {
	io.ReadSeeker
	io.Closer
}

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

// GetTrackInfo returns information about a track
func (s *FileStorageService) GetTrackInfo(ctx context.Context, userID, trackID uint) (*models.Track, error) {
	// Get the track
	var track models.Track
	if err := s.db.First(ctx, &track, trackID); err != nil {
		return nil, errors.New("track not found")
	}

	// Check if the track belongs to a library owned by the user
	var library models.MusicLibrary
	if err := s.db.First(ctx, &library, track.LibraryID); err != nil {
		return nil, errors.New("library not found")
	}

	if library.UserID != userID {
		return nil, errors.New("unauthorized access to track")
	}

	return &track, nil
}

// GetFileStream returns a reader for a track's file
func (s *FileStorageService) GetFileStream(ctx context.Context, userID, trackID uint) (ReadSeekCloser, string, error) {
	// Get the track info
	track, err := s.GetTrackInfo(ctx, userID, trackID)
	if err != nil {
		return nil, "", err
	}

	// Get content type based on file extension
	contentType := s.getContentTypeFromLocation(track.Location)

	// Handle different storage types
	switch track.StorageType {
	case "local":
		// Local file - check if it's a valid local path
		if strings.HasPrefix(track.Location, "file://localhost") {
			return s.getLocalFileStream(track.Location)
		}
		return nil, contentType, errors.New("invalid local file path")
	case "cloud":
		// Cloud file - determine the provider from the location
		if strings.HasPrefix(track.Location, "pcloud://") {
			return s.getPCloudFileStream(track.Location)
		}
		return nil, contentType, errors.New("unsupported cloud provider")
	default:
		return nil, contentType, errors.New("unsupported storage type")
	}
}

// getLocalFileStream returns a reader for a local file
func (s *FileStorageService) getLocalFileStream(location string) (ReadSeekCloser, string, error) {
	// Use the NormalizeTrackLocation function to get a properly formatted path
	path := s.NormalizeTrackLocation(location)

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, "", errors.New("failed to open file: " + err.Error())
	}

	contentType := s.getContentTypeFromLocation(location)
	return file, contentType, nil
}

// getPCloudFileStream returns a reader for a pCloud file
func (s *FileStorageService) getPCloudFileStream(location string) (ReadSeekCloser, string, error) {
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

// NormalizeTrackLocation preprocesses a track location to make it compatible with both Windows and macOS
// This function does not modify the original location in the database, it only normalizes it for frontend use
func (s *FileStorageService) NormalizeTrackLocation(location string) string {
	// If the location is empty, return it as is
	if location == "" {
		return location
	}

	// Handle file:// URLs
	if strings.HasPrefix(location, "file://") {
		// Remove the file://localhost prefix
		location = strings.TrimPrefix(location, "file://localhost")

		// URL-decode the path
		decodedPath, err := url.PathUnescape(location)
		if err == nil {
			location = decodedPath
		}

		// On Windows, if the path starts with a single slash followed by a drive letter and colon,
		// remove the leading slash and convert to Windows path format
		if runtime.GOOS == "windows" && len(location) > 3 && location[0] == '/' && location[2] == ':' {
			// Remove the leading slash
			location = location[1:]
			// Replace forward slashes with backslashes
			location = strings.ReplaceAll(location, "/", "\\")
		}
	}

	// Handle cloud storage URLs (e.g., pcloud://)
	// These are already in a format that doesn't need OS-specific normalization

	return location
}
