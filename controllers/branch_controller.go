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

func ShowBranches(c *gin.Context) {
	var branches []models.Branch

	database.DB.Find(&branches)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List All Branches",
		Data:    branches,
	})
}

func CreateBranch(c *gin.Context) {
	var req structs.BranchCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var branch models.Branch

	// 🔍 cek branch dengan code yang sama (TERMASUK soft deleted)
	err := database.DB.Unscoped().
		Where("code = ?", req.Code).
		First(&branch).Error

	// 👉 jika ditemukan & sudah soft delete → RESTORE
	if err == nil && branch.DeletedAt.Valid {
		branch.Name = req.Name
		branch.DeletedAt = gorm.DeletedAt{}

		if err := database.DB.Unscoped().Save(&branch).Error; err != nil {
			c.JSON(http.StatusBadRequest, structs.ErrorResponse{
				Success: false,
				Message: "Failed to restore branch",
			})
			return
		}

		c.JSON(http.StatusOK, structs.SuccessResponse{
			Success: true,
			Message: "Branch restored successfully",
			Data: structs.BranchResponse{
				ID:        branch.ID,
				Code:      branch.Code,
				Name:      branch.Name,
				CreatedAt: branch.CreatedAt,
				UpdatedAt: branch.UpdatedAt,
			},
		})
		return
	}

	// 👉 jika belum pernah ada → CREATE BARU
	branch = models.Branch{
		Code: req.Code,
		Name: req.Name,
	}

	if err := database.DB.Create(&branch).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create branch",
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Branch created successfully",
		Data: structs.BranchResponse{
			ID:        branch.ID,
			Code:      branch.Code,
			Name:      branch.Name,
			CreatedAt: branch.CreatedAt,
			UpdatedAt: branch.UpdatedAt,
		},
	})
}

func FindBranchById(c *gin.Context) {
	id := c.Param("id")

	var branch models.Branch

	if err := database.DB.First(&branch, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Branch not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Branch found",
		Data: structs.BranchResponse{
			ID:        branch.ID,
			Code:      branch.Code,
			Name:      branch.Name,
			CreatedAt: branch.CreatedAt,
			UpdatedAt: branch.UpdatedAt,
		},
	})
}

func UpdateBranch(c *gin.Context) {
	id := c.Param("id")

	var req structs.BranchUpdateRequest
	var branch models.Branch

	if err := database.DB.First(&branch, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Branch not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	branch.Code = req.Code
	branch.Name = req.Name

	if err := database.DB.Save(&branch).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update branch",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Branch updated successfully",
		Data: structs.BranchResponse{
			ID:        branch.ID,
			Code:      branch.Code,
			Name:      branch.Name,
			CreatedAt: branch.CreatedAt,
			UpdatedAt: branch.UpdatedAt,
		},
	})
}

func DeleteBranch(c *gin.Context) {
	id := c.Param("id")

	var branch models.Branch

	if err := database.DB.First(&branch, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Branch not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Delete(&branch).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete branch",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Branch deleted successfully",
	})
}
