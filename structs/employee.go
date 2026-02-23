package structs

import "time"

type EmployeeResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	BranchID     uint      `json:"branch_id"`
	DivisionID   uint      `json:"division_id"`
	DepartmentID uint      `json:"department_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type EmployeeCreateRequest struct {
	Name         string `json:"name" binding:"required"`
	DepartmentID uint   `json:"department_id" binding:"required"`
}

type EmployeeUpdateRequest struct {
	Name         string `json:"name" binding:"required"`
	DepartmentID uint   `json:"department_id" binding:"required"`
}
