package controllers

import (
	"backend-golang-api/database"
	"backend-golang-api/helpers"
	"backend-golang-api/models"
	"backend-golang-api/structs"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ShowTickets(c *gin.Context) {
	var tickets []models.Ticket

	err := database.DB.
		Preload("Employee").
		Preload("User").
		Order("created_at DESC").
		Find(&tickets).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch tickets",
		})
		return
	}

	// mapping ke response
	var response []structs.TicketListResponse
	for _, t := range tickets {
		response = append(response, structs.TicketListResponse{
			ID:           t.ID,
			TicketNumber: t.TicketNumber,
			Title:        t.Title,
			Priority:     string(t.Priority),
			Status:       string(t.Status),
			EmployeeName: t.Employee.Name,
			CreatedBy:    t.User.Name,
		})
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List All Tickets",
		Data:    response,
	})
}

func ShowActiveTickets(c *gin.Context) {
	var tickets []models.Ticket

	err := database.DB.
		Preload("Employee").
		Where("status != ?", "CLOSED").
		Order("created_at DESC").
		Find(&tickets).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch active tickets",
		})
		return
	}

	var response []structs.TicketListResponse
	for _, t := range tickets {
		response = append(response, structs.TicketListResponse{
			ID:           t.ID,
			TicketNumber: t.TicketNumber,
			Title:        t.Title,
			Priority:     string(t.Priority),
			Status:       string(t.Status),
			EmployeeName: t.Employee.Name,
		})
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List Active Tickets",
		Data:    response,
	})
}

func CreateTicket(c *gin.Context) {
	var req structs.TicketCreateRequest

	// 🔎 Validasi request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// 🔍 Ambil data employee
	var employee models.Employee
	if err := database.DB.First(&employee, req.EmployeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Employee not found",
		})
		return
	}

	var createdTicket models.Ticket
	userID := c.GetUint("user_id")

	// 🔥 Gunakan Transaction supaya atomic
	err := database.DB.Transaction(func(tx *gorm.DB) error {

		// 1️⃣ Insert dulu TANPA ticket_number
		ticket := models.Ticket{
			Title:       req.Title,
			Description: req.Description,
			Priority:    models.TicketPriority(req.Priority),
			Status:      "OPEN",

			EmployeeID:   employee.ID,
			DepartmentID: employee.DepartmentID,
			DivisionID:   employee.DivisionID,
			BranchID:     employee.BranchID,
			UserID:       &userID,
		}

		if err := tx.Create(&ticket).Error; err != nil {
			return err
		}

		// 2️⃣ Generate ticket_number dari ID (SAFE)
		ticketNumber := fmt.Sprintf("IT-%05d", ticket.ID)

		if err := tx.Model(&ticket).
			Update("ticket_number", ticketNumber).Error; err != nil {
			return err
		}

		createdTicket = ticket
		createdTicket.TicketNumber = ticketNumber

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create ticket",
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Ticket created successfully",
		Data:    createdTicket,
	})
}

func ShowTicketDetail(c *gin.Context) {
	ticketNumber := c.Param("ticket_number")

	var ticket models.Ticket

	err := database.DB.
		Preload("Employee").
		Preload("User").
		Preload("Branch").
		Preload("Division").
		Preload("Department").
		Where("ticket_number = ?", ticketNumber).
		First(&ticket).Error

	if err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Ticket not found",
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Ticket detail",
		Data:    ticket,
	})
}

func UpdateTicket(c *gin.Context) {
	ticketNumber := c.Param("ticket_number")

	var req structs.TicketUpdateRequest

	// 🔎 Validasi
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var ticket models.Ticket

	// 🔍 Cek ticket
	if err := database.DB.
		Where("ticket_number = ?", ticketNumber).
		First(&ticket).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Ticket not found",
		})
		return
	}

	// 🔄 Update field
	ticket.Title = req.Title
	ticket.Description = req.Description
	ticket.Priority = models.TicketPriority(req.Priority)
	ticket.Status = models.TicketStatus(req.Status)
	ticket.EmployeeID = req.EmployeeID

	if err := database.DB.Save(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update ticket",
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Ticket updated successfully",
		Data:    ticket,
	})
}

func DeleteTicket(c *gin.Context) {
	ticketNumber := c.Param("ticket_number")

	var ticket models.Ticket

	// 🔍 Cek ticket
	if err := database.DB.
		Where("ticket_number = ?", ticketNumber).
		First(&ticket).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Ticket not found",
		})
		return
	}

	// 🗑 Delete
	if err := database.DB.Delete(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete ticket",
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Ticket deleted successfully",
	})
}
