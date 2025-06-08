package services

import (
	"context"
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/dinis/musync/internal/database"
	"github.com/dinis/musync/internal/models"
)

// XML structures for parsing Rekordbox XML
type RekordboxXML struct {
	XMLName    xml.Name            `xml:"DJ_PLAYLISTS"`
	Version    string              `xml:"Version,attr"`
	Product    RekordboxProduct    `xml:"PRODUCT"`
	Collection RekordboxCollection `xml:"COLLECTION"`
	Playlists  RekordboxPlaylists  `xml:"PLAYLISTS"`
}

type RekordboxProduct struct {
	Name    string `xml:"Name,attr"`
	Version string `xml:"Version,attr"`
	Company string `xml:"Company,attr"`
}

type RekordboxCollection struct {
	Entries int              `xml:"Entries,attr"`
	Tracks  []RekordboxTrack `xml:"TRACK"`
}

type RekordboxTrack struct {
	TrackID     string           `xml:"TrackID,attr"`
	Name        string           `xml:"Name,attr"`
	Artist      string           `xml:"Artist,attr"`
	Composer    string           `xml:"Composer,attr"`
	Album       string           `xml:"Album,attr"`
	Grouping    string           `xml:"Grouping,attr"`
	Genre       string           `xml:"Genre,attr"`
	Kind        string           `xml:"Kind,attr"`
	Size        string           `xml:"Size,attr"`
	TotalTime   string           `xml:"TotalTime,attr"`
	DiscNumber  string           `xml:"DiscNumber,attr"`
	TrackNumber string           `xml:"TrackNumber,attr"`
	Year        string           `xml:"Year,attr"`
	AverageBpm  string           `xml:"AverageBpm,attr"`
	DateAdded   string           `xml:"DateAdded,attr"`
	BitRate     string           `xml:"BitRate,attr"`
	SampleRate  string           `xml:"SampleRate,attr"`
	Comments    string           `xml:"Comments,attr"`
	PlayCount   string           `xml:"PlayCount,attr"`
	Rating      string           `xml:"Rating,attr"`
	Location    string           `xml:"Location,attr"`
	Remixer     string           `xml:"Remixer,attr"`
	Tonality    string           `xml:"Tonality,attr"`
	Label       string           `xml:"Label,attr"`
	Mix         string           `xml:"Mix,attr"`
	Tempo       []RekordboxTempo `xml:"TEMPO"`
}

type RekordboxTempo struct {
	Inizio  string `xml:"Inizio,attr"`
	Bpm     string `xml:"Bpm,attr"`
	Metro   string `xml:"Metro,attr"`
	Battito string `xml:"Battito,attr"`
}

type RekordboxPlaylists struct {
	Nodes []RekordboxNode `xml:"NODE"`
}

type RekordboxNode struct {
	Type    string                   `xml:"Type,attr"`
	Name    string                   `xml:"Name,attr"`
	Count   string                   `xml:"Count,attr"`
	KeyType string                   `xml:"KeyType,attr"`
	Entries string                   `xml:"Entries,attr"`
	Nodes   []RekordboxNode          `xml:"NODE"`
	Tracks  []RekordboxPlaylistTrack `xml:"TRACK"`
}

type RekordboxPlaylistTrack struct {
	Key string `xml:"Key,attr"`
}

// MusicLibraryService handles operations related to music libraries
type MusicLibraryService struct {
	db *database.DB
}

// NewMusicLibraryService creates a new MusicLibraryService
func NewMusicLibraryService(db *database.DB) *MusicLibraryService {
	return &MusicLibraryService{
		db: db,
	}
}

// UploadLibrary parses and stores a music library from an XML file
func (s *MusicLibraryService) UploadLibrary(ctx context.Context, userID uint, name string, xmlReader io.Reader) (uint, error) {
	// Parse the XML file
	var rekordboxXML RekordboxXML
	if err := xml.NewDecoder(xmlReader).Decode(&rekordboxXML); err != nil {
		return 0, errors.New("failed to parse XML file")
	}

	// Begin a transaction
	var libraryID uint
	err := s.db.Transaction(ctx, func(tx *database.DB) error {
		// Create the music library
		library := models.MusicLibrary{
			UserID:      userID,
			Name:        name,
			Source:      "rekordbox",
			Version:     rekordboxXML.Version,
			ProductName: rekordboxXML.Product.Name,
			Company:     rekordboxXML.Product.Company,
		}

		if err := tx.Create(ctx, &library); err != nil {
			return err
		}

		libraryID = library.ID

		// Process tracks
		trackIDMap := make(map[string]uint) // Map original TrackID to our Track ID
		for _, rbTrack := range rekordboxXML.Collection.Tracks {
			track, err := s.convertRekordboxTrack(ctx, rbTrack, library.ID)
			if err != nil {
				return err
			}

			if err := tx.Create(ctx, &track); err != nil {
				return err
			}

			trackIDMap[rbTrack.TrackID] = track.ID

			// Process tempo markers
			for _, rbTempo := range rbTrack.Tempo {
				tempo, err := s.convertRekordboxTempo(ctx, rbTempo, track.ID)
				if err != nil {
					return err
				}

				if err := tx.Create(ctx, &tempo); err != nil {
					return err
				}
			}
		}

		// Process playlists
		for _, rbNode := range rekordboxXML.Playlists.Nodes {
			if err := s.processRekordboxNode(ctx, tx, rbNode, library.ID, nil, trackIDMap); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return libraryID, nil
}

// GetLibraries returns all music libraries for a user
func (s *MusicLibraryService) GetLibraries(ctx context.Context, userID uint) ([]models.MusicLibrary, error) {
	var libraries []models.MusicLibrary
	if err := s.db.Where(ctx, "user_id = ?", userID).Find(ctx, &libraries); err != nil {
		return nil, err
	}
	return libraries, nil
}

// GetLibrary returns a specific music library
func (s *MusicLibraryService) GetLibrary(ctx context.Context, userID, libraryID uint) (*models.MusicLibrary, error) {
	var library models.MusicLibrary
	if err := s.db.Where(ctx, "id = ? AND user_id = ?", libraryID, userID).First(ctx, &library); err != nil {
		return nil, err
	}
	return &library, nil
}

// GetTracks returns all tracks in a library
func (s *MusicLibraryService) GetTracks(ctx context.Context, userID, libraryID uint) ([]models.Track, error) {
	// First check if the library belongs to the user
	if _, err := s.GetLibrary(ctx, userID, libraryID); err != nil {
		return nil, err
	}

	var tracks []models.Track
	if err := s.db.Where(ctx, "library_id = ?", libraryID).Find(ctx, &tracks); err != nil {
		return nil, err
	}
	return tracks, nil
}

// GetPlaylists returns all playlists in a library
func (s *MusicLibraryService) GetPlaylists(ctx context.Context, userID, libraryID uint) ([]models.Playlist, error) {
	// First check if the library belongs to the user
	if _, err := s.GetLibrary(ctx, userID, libraryID); err != nil {
		return nil, err
	}

	var playlists []models.Playlist
	if err := s.db.Where(ctx, "library_id = ?", libraryID).Find(ctx, &playlists); err != nil {
		return nil, err
	}
	return playlists, nil
}

// GetPlaylistTracks returns all tracks in a playlist
func (s *MusicLibraryService) GetPlaylistTracks(ctx context.Context, userID, playlistID uint) ([]models.Track, error) {
	// First get the playlist
	var playlist models.Playlist
	if err := s.db.First(ctx, &playlist, playlistID); err != nil {
		return nil, err
	}

	// Check if the library belongs to the user
	if _, err := s.GetLibrary(ctx, userID, playlist.LibraryID); err != nil {
		return nil, err
	}

	// Get playlist tracks
	var playlistTracks []models.PlaylistTrack
	if err := s.db.Where(ctx, "playlist_id = ?", playlistID).Find(ctx, &playlistTracks); err != nil {
		return nil, err
	}

	// Get the actual tracks
	var tracks []models.Track
	for _, pt := range playlistTracks {
		var track models.Track
		if err := s.db.Where(ctx, "library_id = ? AND track_id = ?", playlist.LibraryID, pt.TrackKey).First(ctx, &track); err != nil {
			continue // Skip tracks that can't be found
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

// Helper functions

// convertRekordboxTrack converts a RekordboxTrack to a Track model
func (s *MusicLibraryService) convertRekordboxTrack(ctx context.Context, rbTrack RekordboxTrack, libraryID uint) (models.Track, error) {
	size, _ := strconv.ParseInt(rbTrack.Size, 10, 64)
	totalTime, _ := strconv.Atoi(rbTrack.TotalTime)
	discNumber, _ := strconv.Atoi(rbTrack.DiscNumber)
	trackNumber, _ := strconv.Atoi(rbTrack.TrackNumber)
	year, _ := strconv.Atoi(rbTrack.Year)
	averageBpm, _ := strconv.ParseFloat(rbTrack.AverageBpm, 64)
	bitRate, _ := strconv.Atoi(rbTrack.BitRate)
	sampleRate, _ := strconv.Atoi(rbTrack.SampleRate)
	playCount, _ := strconv.Atoi(rbTrack.PlayCount)
	rating, _ := strconv.Atoi(rbTrack.Rating)

	dateAdded, _ := time.Parse("2006-01-02", rbTrack.DateAdded)

	return models.Track{
		LibraryID:   libraryID,
		TrackID:     rbTrack.TrackID,
		Name:        rbTrack.Name,
		Artist:      rbTrack.Artist,
		Composer:    rbTrack.Composer,
		Album:       rbTrack.Album,
		Grouping:    rbTrack.Grouping,
		Genre:       rbTrack.Genre,
		Kind:        rbTrack.Kind,
		Size:        size,
		TotalTime:   totalTime,
		DiscNumber:  discNumber,
		TrackNumber: trackNumber,
		Year:        year,
		AverageBpm:  averageBpm,
		DateAdded:   dateAdded,
		BitRate:     bitRate,
		SampleRate:  sampleRate,
		Comments:    rbTrack.Comments,
		PlayCount:   playCount,
		Rating:      rating,
		Location:    rbTrack.Location,
		Remixer:     rbTrack.Remixer,
		Tonality:    rbTrack.Tonality,
		Label:       rbTrack.Label,
		Mix:         rbTrack.Mix,
	}, nil
}

// convertRekordboxTempo converts a RekordboxTempo to a Tempo model
func (s *MusicLibraryService) convertRekordboxTempo(ctx context.Context, rbTempo RekordboxTempo, trackID uint) (models.Tempo, error) {
	inizio, _ := strconv.ParseFloat(rbTempo.Inizio, 64)
	bpm, _ := strconv.ParseFloat(rbTempo.Bpm, 64)
	battito, _ := strconv.Atoi(rbTempo.Battito)

	return models.Tempo{
		TrackID: trackID,
		Inizio:  inizio,
		Bpm:     bpm,
		Metro:   rbTempo.Metro,
		Battito: battito,
	}, nil
}

// processRekordboxNode recursively processes a RekordboxNode (playlist or folder)
func (s *MusicLibraryService) processRekordboxNode(ctx context.Context, tx *database.DB, rbNode RekordboxNode, libraryID uint, parentID *uint, trackIDMap map[string]uint) error {
	// Determine if this is a folder or playlist
	nodeType, _ := strconv.Atoi(rbNode.Type)

	// Create the playlist/folder
	playlist := models.Playlist{
		LibraryID: libraryID,
		Name:      rbNode.Name,
		Type:      nodeType,
		ParentID:  parentID,
	}

	if err := tx.Create(ctx, &playlist); err != nil {
		return err
	}

	// Process tracks if this is a playlist
	if nodeType == 1 && len(rbNode.Tracks) > 0 {
		for _, rbTrack := range rbNode.Tracks {
			playlistTrack := models.PlaylistTrack{
				PlaylistID: playlist.ID,
				TrackKey:   rbTrack.Key,
			}

			if err := tx.Create(ctx, &playlistTrack); err != nil {
				return err
			}
		}
	}

	// Process child nodes recursively
	for _, childNode := range rbNode.Nodes {
		if err := s.processRekordboxNode(ctx, tx, childNode, libraryID, &playlist.ID, trackIDMap); err != nil {
			return err
		}
	}

	return nil
}
