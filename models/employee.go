package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Name         string `json:"name" gorm:"type:varchar(100);not null"`
	BranchID     uint   `json:"branch_id" gorm:"not null;index"`
	DivisionID   uint   `json:"division_id" gorm:"not null;index"`
	DepartmentID uint   `json:"department_id" gorm:"not null;index"`

	Branch     Branch     `json:"branch" gorm:"foreignKey:BranchID"`
	Division   Division   `json:"division" gorm:"foreignKey:DivisionID"`
	Department Department `json:"department" gorm:"foreignKey:DepartmentID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
