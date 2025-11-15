package model

import (
	"time"

	"gorm.io/gorm"
)

type UserSession struct {
	gorm.Model
	UserID             uint      `gorm:"not null;index"`
	User               User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	HashedRefreshToken string    `gorm:"uniqueIndex;not null"`
	UserAgent          string    `gorm:"size:255"`
	ClientIP           string    `gorm:"size:45"`
	IsRevoked          bool      `gorm:"not null;default:false"`
	ExpiresAt          time.Time `gorm:"not null"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
