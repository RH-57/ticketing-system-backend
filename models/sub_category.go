package models

import (
	"time"

	"gorm.io/gorm"
)

type SubCategory struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Name       string `json:"name" gorm:"type:varchar(100);not null"`
	Slug       string `json:"slug" gorm:"type:varchar(120);uniqueIndex;not null"`
	CategoryID uint   `json:"category_id" gorm:"not null;index"`

	Category Category `json:"category" gorm:"foreignKey:CategoryID"`
	Items    []Item   `json:"sub_sub_categories" gorm:"foreignKey:SubCategoryID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
