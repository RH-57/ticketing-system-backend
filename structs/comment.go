package structs

import "time"

type CommentResponse struct {
	ID uint `json:"id"`

	TicketID uint `json:"ticket_id"`
	UserID   uint `json:"user_id"`

	CategoryID    uint `json:"category_id"`
	SubCategoryID uint `json:"sub_category_id"`
	ItemID        uint `json:"item_id"`

	Type        string `json:"type"`
	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`
}

type CommentCreateRequest struct {
	TicketID uint `json:"ticket_id" binding:"required"`

	CategoryID    uint `json:"category_id" binding:"required"`
	SubCategoryID uint `json:"sub_category_id" binding:"required"`
	ItemID        uint `json:"item_id" binding:"required"`

	Type string `json:"type" binding:"required,oneof=Malfunction Human_Error Install Other"`

	Description string `json:"description"`
}

type CommentUpdateRequest struct {
	CategoryID    uint `json:"category_id" binding:"required"`
	SubCategoryID uint `json:"sub_category_id" binding:"required"`
	ItemID        uint `json:"item_id" binding:"required"`

	Type string `json:"type" binding:"required,oneof=Malfunction Human_Error Install Other"`

	Description string `json:"description"`
}
