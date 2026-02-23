package models

import (
	"time"
)

type CommentType string

const (
	TypeMalfunction CommentType = "Malfunction"
	TypeHumanError  CommentType = "Human_Error"
	TypeInstall     CommentType = "Install"
	TypeOther       CommentType = "Other"
)

type Comment struct {
	ID uint `json:"id" gorm:"primaryKey"`

	// ======================
	// Foreign Keys
	// ======================

	TicketID uint `json:"ticket_id" gorm:"index;not null"`
	UserID   uint `json:"user_id" gorm:"index;not null"`

	CategoryID    uint `json:"category_id" gorm:"index;not null"`
	SubCategoryID uint `json:"sub_category_id" gorm:"index;not null"`
	ItemID        uint `json:"item_id" gorm:"index;not null"`

	// ======================
	// Relations
	// ======================

	Ticket      Ticket      `json:"ticket" gorm:"foreignKey:TicketID"`
	User        User        `json:"user" gorm:"foreignKey:UserID"`
	Category    Category    `json:"category" gorm:"foreignKey:CategoryID"`
	SubCategory SubCategory `json:"sub_category" gorm:"foreignKey:SubCategoryID"`
	Item        Item        `json:"item" gorm:"foreignKey:ItemID"`

	// ======================
	// Enum
	// ======================

	Type CommentType `json:"type" gorm:"type:varchar(50);default:'Malfunction'"`

	Description string `json:"description" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at"`
}
