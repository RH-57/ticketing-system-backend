package structs

import "time"

type EmployeeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`

	Branch struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"branch"`

	Division struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"division"`

	Department struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"department"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EmployeeCreateRequest struct {
	Name         string `json:"name" binding:"required"`
	DepartmentID uint   `json:"department_id" binding:"required"`
}

type EmployeeUpdateRequest struct {
	Name         string `json:"name" binding:"required"`
	DepartmentID uint   `json:"department_id" binding:"required"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	LastPage    int   `json:"last_page"`
}

type PaginatedResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}
