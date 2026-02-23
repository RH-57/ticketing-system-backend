package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"type:varchar(100);not null"`
	Slug string `json:"slug" gorm:"type:varchar(120);uniqueIndex;not null"`

	SubCategories []SubCategory `json:"sub_categories" gorm:"foreignKey:CategoryID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;uniqueIndex:idx_code_deleted"`
}
