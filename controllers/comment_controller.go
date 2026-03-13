package controllers

import (
	"backend-golang-api/database"
	"backend-golang-api/helpers"
	"backend-golang-api/models"
	"backend-golang-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCommentAndCloseTicket(c *gin.Context) {
	// request
	var req structs.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// 🔍 cek ticket
	var ticket models.Ticket
	if err := database.DB.First(&ticket, req.TicketID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Ticket not found",
		})
		return
	}

	// Cek apakah tiket sudah closed sebelumnya (opsional, tapi bagus untuk validasi)
	if ticket.Status == "CLOSED" {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Ticket is already closed",
		})
		return
	}

	// Ambil user ID dari JWT/Auth Middleware
	userID := c.GetUint("user_id")

	var comment models.Comment

	// 🔥 Gunakan Transaction agar Insert Comment & Update Status bersifat Atomic
	err := database.DB.Transaction(func(tx *gorm.DB) error {

		// 👉 CREATE BARU comment
		comment = models.Comment{
			TicketID:      ticket.ID,
			UserID:        userID,
			CategoryID:    req.CategoryID,
			SubCategoryID: req.SubCategoryID,
			ItemID:        req.ItemID,
			Type:          models.CommentType(req.Type),
			Description:   req.Description,
		}

		if err := tx.Create(&comment).Error; err != nil {
			return err
		}

		// 👉 UPDATE status ticket -> CLOSED
		if err := tx.Model(&ticket).Update("status", "CLOSED").Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to submit comment and close ticket",
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Get All Departments",
		Data:    comment,
	})
}

// Opsional: Fungsi untuk melihat comment di detail tiket nantinya
func ShowCommentByTicket(c *gin.Context) {
	ticketID := c.Param("ticketId")

	// 🔍 cek ticket
	var ticket models.Ticket
	if err := database.DB.First(&ticket, ticketID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Ticket not found",
		})
		return
	}

	var comment models.Comment
	// 🔍 ambil comment (karena 1 tiket = 1 form penyelesaian)
	if err := database.DB.Where("ticket_id = ?", ticket.ID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Comment not found for this ticket",
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Get All Departments",
		Data:    comment,
	})
}
