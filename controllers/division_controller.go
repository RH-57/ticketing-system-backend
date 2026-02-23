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

func ShowDivisionsByBranch(c *gin.Context) {
	branchID := c.Param("id")

	var branch models.Branch

	if err := database.DB.First(&branch, branchID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Branch not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var divisions []models.Division

	if err := database.DB.
		Where("branch_id = ?", branchID).
		Order("created_at ASC").
		Find(&divisions).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch divisions",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List all divisions",
		Data:    divisions,
	})
}

func CreateDivision(c *gin.Context) {
	branchID := c.Param("id")

	// 🔍 cek branch
	var branch models.Branch
	if err := database.DB.First(&branch, branchID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Branch not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// request
	var req structs.DivisionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var division models.Division

	// 🔍 cek division dengan code + branch_id (TERMASUK soft deleted)
	err := database.DB.Unscoped().
		Where("code = ? AND branch_id = ?", req.Code, branch.ID).
		First(&division).Error

	// 👉 jika ditemukan & sudah soft delete → RESTORE
	if err == nil && division.DeletedAt.Valid {
		division.Name = req.Name
		division.DeletedAt = gorm.DeletedAt{}

		if err := database.DB.Unscoped().Save(&division).Error; err != nil {
			c.JSON(http.StatusBadRequest, structs.ErrorResponse{
				Success: false,
				Message: "Failed to restore division",
			})
			return
		}

		c.JSON(http.StatusOK, structs.SuccessResponse{
			Success: true,
			Message: "Division restored successfully",
			Data: structs.DivisionResponse{
				ID:        division.ID,
				Code:      division.Code,
				Name:      division.Name,
				BranchID:  division.BranchID,
				CreatedAt: division.CreatedAt,
				UpdatedAt: division.UpdatedAt,
			},
		})
		return
	}

	// 👉 jika belum pernah ada → CREATE BARU
	division = models.Division{
		Code:     req.Code,
		Name:     req.Name,
		BranchID: branch.ID,
	}

	if err := database.DB.Create(&division).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create division",
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Division created successfully",
		Data: structs.DivisionResponse{
			ID:        division.ID,
			Code:      division.Code,
			Name:      division.Name,
			BranchID:  division.BranchID,
			CreatedAt: division.CreatedAt,
			UpdatedAt: division.UpdatedAt,
		},
	})
}

func UpdateDivision(c *gin.Context) {
	branchID := c.Param("id")
	divisionID := c.Param("divisionId")

	// 🔍 cek branch
	var branch models.Branch
	if err := database.DB.First(&branch, branchID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Branch not found",
		})
		return
	}

	// 🔍 cek division (harus milik branch tsb)
	var division models.Division
	if err := database.DB.
		Where("id = ? AND branch_id = ?", divisionID, branch.ID).
		First(&division).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Division not found",
		})
		return
	}

	// request
	var req structs.DivisionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// 🔍 cek konflik code di branch yang sama (kecuali diri sendiri)
	var conflict models.Division
	err := database.DB.
		Where("code = ? AND branch_id = ? AND id <> ?", req.Code, branch.ID, division.ID).
		First(&conflict).Error

	if err == nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Division code already exists in this branch",
		})
		return
	}

	// update
	division.Code = req.Code
	division.Name = req.Name

	if err := database.DB.Save(&division).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update division",
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Division updated successfully",
		Data: structs.DivisionResponse{
			ID:        division.ID,
			Code:      division.Code,
			Name:      division.Name,
			BranchID:  division.BranchID,
			CreatedAt: division.CreatedAt,
			UpdatedAt: division.UpdatedAt,
		},
	})
}

func DeleteDivision(c *gin.Context) {
	branchID := c.Param("id")
	divisionID := c.Param("divisionId")

	// 🔍 cek division + branch
	var division models.Division
	if err := database.DB.
		Where("id = ? AND branch_id = ?", divisionID, branchID).
		First(&division).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Division not found",
		})
		return
	}

	// soft delete
	if err := database.DB.Delete(&division).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete division",
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Division deleted successfully",
	})
}
