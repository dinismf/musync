package models

import (
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	UserID      uint `gorm:"uniqueIndex;not null"`
	DisplayName string
	Bio         string
	AvatarURL   string
	Location    string
	Website     string
	SocialLinks []SocialLink
}
