package models

import "time"

const (
	RoleSuperuser = "superuser"
	RoleAdmin     = "admin"
	RoleView      = "view"
	Technician    = "technician"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"type:varchar(100);not null"`
	Email    string `json:"email" gorm:"type:varchar(100);unique;not null"`
	Username string `json:"username" gorm:"type:varchar(100);unique;not null"`
	Password string `json:"-" gorm:"not null"`
	Role     string `json:"role" gorm:"type:varchar(20);not null;default:'admin'"`
	IsActive bool   `json:"is_active" gorm:"default:true"`
	Avatar   string `json:"avatar" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
