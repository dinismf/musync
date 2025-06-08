package models

import (
	"gorm.io/gorm"
)

type Release struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Type        string `gorm:"not null"` // album, single, ep
	ReleaseDate string
	ImageURL    string
	SpotifyID   string `gorm:"uniqueIndex"`
	ArtistID    uint   `gorm:"not null"`
	Artist      Artist
	LabelID     uint
	Label       Label
	FeedItems   []FeedItem
}
