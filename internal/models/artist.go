package models

import (
	"gorm.io/gorm"
)

type Artist struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	ImageURL    string
	SpotifyID   string `gorm:"uniqueIndex"`
	Releases    []Release
	Labels      []Label `gorm:"many2many:artist_labels;"`
	Followers   []User  `gorm:"many2many:user_following_artists;"`
}
