package structs

import "time"

type ItemCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type ItemUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}

type ItemResponse struct {
	ID            uint                `json:"id"`
	Name          string              `json:"name"`
	Slug          string              `json:"slug"`
	SubCategoryID uint                `json:"sub_category_id"`
	SubCategory   SubCategoryResponse `json:"sub_category"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}
