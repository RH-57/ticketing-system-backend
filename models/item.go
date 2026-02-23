package models

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Name          string `json:"name" gorm:"type:varchar(100);not null"`
	Slug          string `json:"slug" gorm:"type:varchar(120);uniqueIndex;not null"`
	SubCategoryID uint   `json:"sub_category_id" gorm:"not null;index"`

	SubCategory SubCategory `json:"sub_category" gorm:"foreignKey:SubCategoryID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
