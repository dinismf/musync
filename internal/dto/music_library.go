package dto

import (
	"time"

	"github.com/dinis/musync/internal/models"
)

// LibraryResponse represents the response structure for a music library
type LibraryResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Name      string    `json:"name"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TrackResponse represents the response structure for a track
type TrackResponse struct {
	ID          uint      `json:"id"`
	LibraryID   uint      `json:"library_id"`
	Name        string    `json:"title"`
	Artist      string    `json:"artist"`
	Album       string    `json:"album"`
	Genre       string    `json:"genre"`
	TotalTime   int       `json:"duration"`
	Year        int       `json:"year"`
	Location    string    `json:"location,omitempty"`
	StorageType string    `json:"storage_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PlaylistResponse represents the response structure for a playlist
type PlaylistResponse struct {
	ID        uint      `json:"id"`
	LibraryID uint      `json:"library_id"`
	Name      string    `json:"name"`
	Type      int       `json:"type"`
	ParentID  *uint     `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToLibraryResponse converts a MusicLibrary model to a LibraryResponse DTO
func ToLibraryResponse(library models.MusicLibrary) LibraryResponse {
	return LibraryResponse{
		ID:        library.ID,
		UserID:    library.UserID,
		Name:      library.Name,
		Source:    library.Source,
		CreatedAt: library.CreatedAt,
		UpdatedAt: library.UpdatedAt,
	}
}

// ToLibraryResponses converts a slice of MusicLibrary models to a slice of LibraryResponse DTOs
func ToLibraryResponses(libraries []models.MusicLibrary) []LibraryResponse {
	responses := make([]LibraryResponse, len(libraries))
	for i, library := range libraries {
		responses[i] = ToLibraryResponse(library)
	}
	return responses
}

// ToTrackResponse converts a Track model to a TrackResponse DTO
func ToTrackResponse(track models.Track) TrackResponse {
	return TrackResponse{
		ID:          track.ID,
		LibraryID:   track.LibraryID,
		Name:        track.Name,
		Artist:      track.Artist,
		Album:       track.Album,
		Genre:       track.Genre,
		TotalTime:   track.TotalTime,
		Year:        track.Year,
		Location:    track.Location,
		StorageType: track.StorageType,
		CreatedAt:   track.CreatedAt,
		UpdatedAt:   track.UpdatedAt,
	}
}

// ToTrackResponses converts a slice of Track models to a slice of TrackResponse DTOs
func ToTrackResponses(tracks []models.Track) []TrackResponse {
	responses := make([]TrackResponse, len(tracks))
	for i, track := range tracks {
		responses[i] = ToTrackResponse(track)
	}
	return responses
}

// ToPlaylistResponse converts a Playlist model to a PlaylistResponse DTO
func ToPlaylistResponse(playlist models.Playlist) PlaylistResponse {
	return PlaylistResponse{
		ID:        playlist.ID,
		LibraryID: playlist.LibraryID,
		Name:      playlist.Name,
		Type:      playlist.Type,
		ParentID:  playlist.ParentID,
		CreatedAt: playlist.CreatedAt,
		UpdatedAt: playlist.UpdatedAt,
	}
}

// ToPlaylistResponses converts a slice of Playlist models to a slice of PlaylistResponse DTOs
func ToPlaylistResponses(playlists []models.Playlist) []PlaylistResponse {
	responses := make([]PlaylistResponse, len(playlists))
	for i, playlist := range playlists {
		responses[i] = ToPlaylistResponse(playlist)
	}
	return responses
}
