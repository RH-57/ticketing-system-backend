package models

import (
	"time"

	"gorm.io/gorm"
)

type Branch struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Code      string         `json:"code" gorm:"type:varchar(10);unique;not null"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"`
	Divisions []Division     `json:"divisions,omitempty" gorm:"foreignKey:BranchID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;uniqueIndex:idx_code_deleted"`
}
