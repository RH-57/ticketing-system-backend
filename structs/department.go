package structs

import "time"

type DepartmentResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	DivisionID uint      `json:"division_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DepartmentCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type DepartmentUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}
