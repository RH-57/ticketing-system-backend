package models

import (
	"time"

	"gorm.io/gorm"
)

type Department struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Name       string `json:"name" gorm:"type:varchar(100);not null"`
	DivisionID uint   `json:"division_id" gorm:"not null;index"`

	Division Division `json:"division"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;uniqueIndex:idx_code_deleted"`
}
