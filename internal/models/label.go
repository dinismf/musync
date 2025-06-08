package models

import (
	"gorm.io/gorm"
)

type Label struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	ImageURL    string
	Website     string
	SpotifyID   string   `gorm:"uniqueIndex"`
	Artists     []Artist `gorm:"many2many:artist_labels;"`
	Releases    []Release
	Followers   []User `gorm:"many2many:user_following_labels;"`
}
