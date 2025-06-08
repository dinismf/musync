package models

import (
	"gorm.io/gorm"
)

type SocialLink struct {
	gorm.Model
	ProfileID uint   `gorm:"not null"`
	Platform  string `gorm:"not null"` // e.g., "spotify", "soundcloud", "instagram"
	URL       string `gorm:"not null"`
}
