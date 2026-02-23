package models

import (
	"time"

	"gorm.io/gorm"
)

type TicketStatus string
type TicketPriority string

const (
	StatusOpen     TicketStatus = "OPEN"
	StatusProcess  TicketStatus = "PROCESS"
	StatusPending  TicketStatus = "PENDING"
	StatusResolved TicketStatus = "RESOLVED"
	StatusClosed   TicketStatus = "CLOSED"

	PriorityLow      TicketPriority = "LOW"
	PriorityMedium   TicketPriority = "MEDIUM"
	PriorityHigh     TicketPriority = "HIGH"
	PriorityCritical TicketPriority = "CRITICAL"
)

type Ticket struct {
	ID uint `json:"id" gorm:"primaryKey"`

	TicketNumber string `json:"ticket_number" gorm:"type:varchar(50);uniqueIndex;not null"`
	Title        string `json:"title" gorm:"type:varchar(255);not null"`
	Description  string `json:"description" gorm:"type:text"`

	// ======================
	// Foreign Keys
	// ======================

	BranchID     uint  `json:"branch_id" gorm:"index;not null"`
	DivisionID   uint  `json:"division_id" gorm:"index;not null"`
	DepartmentID uint  `json:"department_id" gorm:"index;not null"`
	EmployeeID   uint  `json:"employee_id" gorm:"index;not null"`
	UserID       *uint `json:"user_id" gorm:"index;not null"` // reported by

	// ======================
	// Relations
	// ======================

	Branch     Branch     `json:"branch" gorm:"foreignKey:BranchID"`
	Division   Division   `json:"division" gorm:"foreignKey:DivisionID"`
	Department Department `json:"department" gorm:"foreignKey:DepartmentID"`
	Employee   Employee   `json:"employee" gorm:"foreignKey:EmployeeID"`
	User       User       `json:"reported_by" gorm:"foreignKey:UserID"`

	Comments []Comment `json:"comments" gorm:"foreignKey:TicketID"`

	// ======================
	// Enum Fields
	// ======================

	Status   TicketStatus   `json:"status" gorm:"type:varchar(20);default:'OPEN';index"`
	Priority TicketPriority `json:"priority" gorm:"type:varchar(20);default:'MEDIUM';index"`

	ClosedAt *time.Time `json:"closed_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
