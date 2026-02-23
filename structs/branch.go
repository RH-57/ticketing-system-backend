package structs

import "time"

type BranchResponse struct {
	ID        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BranchCreateRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type BranchUpdateRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}
