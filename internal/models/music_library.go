package models

import (
	"time"

	"gorm.io/gorm"
)

// MusicLibrary represents a user's music library from a specific source (e.g., Rekordbox)
type MusicLibrary struct {
	gorm.Model
	UserID      uint   `gorm:"not null"`
	User        User   `gorm:"foreignKey:UserID"`
	Name        string `gorm:"not null"`
	Source      string `gorm:"not null"` // e.g., "rekordbox", "serato", etc.
	Version     string
	ProductName string
	Company     string
	Tracks      []Track    `gorm:"foreignKey:LibraryID"`
	Playlists   []Playlist `gorm:"foreignKey:LibraryID"`
}

// Track represents a music track in a library
type Track struct {
	gorm.Model
	LibraryID   uint   `gorm:"not null"`
	TrackID     string `gorm:"not null"` // Original ID from the source library
	Name        string `gorm:"not null"`
	Artist      string
	Composer    string
	Album       string
	Grouping    string
	Genre       string
	Kind        string // File type (e.g., "WAV File")
	Size        int64  // File size in bytes
	TotalTime   int    // Duration in seconds
	DiscNumber  int
	TrackNumber int
	Year        int
	AverageBpm  float64
	DateAdded   time.Time
	BitRate     int
	SampleRate  int
	Comments    string
	PlayCount   int
	Rating      int
	Location    string `gorm:"not null"`                 // File path or URL
	StorageType string `gorm:"not null;default:'local'"` // Storage type: "local" or "cloud"
	Remixer     string
	Tonality    string // Musical key
	Label       string
	Mix         string
	Tempo       []Tempo `gorm:"foreignKey:TrackID"` // One track can have multiple tempo markers
}

// Tempo represents a tempo marker in a track
type Tempo struct {
	gorm.Model
	TrackID uint    `gorm:"not null"`
	Inizio  float64 // Start position in seconds
	Bpm     float64
	Metro   string // Time signature (e.g., "4/4")
	Battito int    // Beat number
}

// Playlist represents a playlist in a library
type Playlist struct {
	gorm.Model
	LibraryID      uint            `gorm:"not null"`
	Name           string          `gorm:"not null"`
	Type           int             // 0 for folder, 1 for playlist
	ParentID       *uint           // For nested playlists/folders
	Parent         *Playlist       `gorm:"foreignKey:ParentID"`
	PlaylistTracks []PlaylistTrack `gorm:"foreignKey:PlaylistID"`
}

// PlaylistTrack represents a track in a playlist
type PlaylistTrack struct {
	gorm.Model
	PlaylistID uint   `gorm:"not null"`
	TrackKey   string `gorm:"not null"` // References Track.TrackID
}
