package models

import (
	"time"

	"gorm.io/gorm"
)

type Division struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"type:varchar(10);not null;uniqueIndex:uniq_code_branch"`
	Name     string `json:"name" gorm:"type:varchar(100);not null"`
	BranchID uint   `json:"branch_id" gorm:"not null;uniqueIndex:uniq_code_branch"`

	Departments []Department `json:"departments,omitempty" gorm:"foreignKey:DivisionID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
