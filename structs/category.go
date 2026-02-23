package structs

import "time"

type CategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type CategoryUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}
