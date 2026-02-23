package models

import "time"

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"type:text;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Revoked   bool      `gorm:"default:false"`
	UserAgent string    `gorm:"type:varchar(255)"`
	IPAddress string    `gorm:"type:varchar(100)"`

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`

	CreatedAt time.Time
}
