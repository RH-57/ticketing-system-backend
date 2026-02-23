package structs

import "time"

type SubCategoryCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type SubCategoryUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}

type SubCategoryResponse struct {
	ID         uint             `json:"id"`
	Name       string           `json:"name"`
	Slug       string           `json:"slug"`
	CategoryID uint             `json:"category_id"`
	Category   CategoryResponse `json:"category"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}
