package models

import (
	"gorm.io/gorm"
)

type FeedItem struct {
	gorm.Model
	UserID    uint `gorm:"not null"`
	User      User
	Type      string `gorm:"not null"` // new_release, artist_update, label_update
	Content   string
	ReleaseID *uint
	Release   *Release
	Read      bool `gorm:"default:false"`
}
