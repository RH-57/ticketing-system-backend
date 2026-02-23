package structs

import "time"

type DivisionResponse struct {
	ID        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	BranchID  uint      `json:"branch_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DivisionCreateRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type DivisionUpdateRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}
