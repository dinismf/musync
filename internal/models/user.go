package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email            string `gorm:"uniqueIndex;not null"`
	PasswordHash     string
	Username         string `gorm:"uniqueIndex;not null"`
	IsEmailVerified  bool   `gorm:"default:false"`
	IsPasswordSet    bool   `gorm:"default:false"`
	VerificationCode string
	ResetCode        string
	ResetExpiresAt   *time.Time
	Profile          Profile
	Feed             []FeedItem
	Following        []Artist `gorm:"many2many:user_following_artists;"`
	Labels           []Label  `gorm:"many2many:user_following_labels;"`
}
