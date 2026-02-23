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

func ShowSubCategoriesByCategory(c *gin.Context) {
	categoryID := c.Param("id")

	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var subCategory []models.SubCategory

	if err := database.DB.
		Preload("Category"). // 🔥 TAMBAHKAN INI
		Where("category_id = ?", categoryID).
		Order("name ASC").
		Find(&subCategory).Error; err != nil {

		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch Sub Category",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Show all sub category",
		Data:    subCategory,
	})
}

func CreateSubCategory(c *gin.Context) {
	categoryID := c.Param("id")

	// cek category
	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// bind request
	var req structs.SubCategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	slug := helpers.GenerateSlug(req.Name)

	var subCategory models.SubCategory

	// cek existing termasuk soft delete
	err := database.DB.Unscoped().
		Where("slug = ? AND category_id = ?", slug, category.ID).
		First(&subCategory).Error

	// ================= RESTORE =================
	if err == nil && subCategory.DeletedAt.Valid {
		subCategory.Name = req.Name
		subCategory.Slug = slug
		subCategory.DeletedAt = gorm.DeletedAt{}

		if err := database.DB.Unscoped().Save(&subCategory).Error; err != nil {
			c.JSON(http.StatusBadRequest, structs.ErrorResponse{
				Success: false,
				Message: "Failed to restore sub category",
			})
			return
		}
	} else {
		// ================= CREATE BARU =================
		subCategory = models.SubCategory{
			Name:       req.Name,
			Slug:       slug,
			CategoryID: category.ID,
		}

		if err := database.DB.Create(&subCategory).Error; err != nil {
			c.JSON(http.StatusBadRequest, structs.ErrorResponse{
				Success: false,
				Message: "Failed to create sub category",
			})
			return
		}
	}

	// 🔥 reload agar relasi category terisi
	if err := database.DB.
		Preload("Category").
		First(&subCategory, subCategory.ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to load sub category",
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Sub Category saved successfully",
		Data:    subCategory,
	})
}

func ShowSubCategoryByID(c *gin.Context) {
	categoryID := c.Param("id")
	subCategoryID := c.Param("subCategoryId")

	var subCategory models.SubCategory

	if err := database.DB.
		Preload("Category").
		Where("id = ? AND category_id = ?", subCategoryID, categoryID).
		First(&subCategory).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Sub Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Sub Category detail",
		Data:    subCategory,
	})
}

func UpdateSubCategory(c *gin.Context) {
	categoryID := c.Param("id")
	subCategoryID := c.Param("subCategoryId")

	// cek category
	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Category not found",
		})
		return
	}

	// cek subcategory
	var subCategory models.SubCategory
	if err := database.DB.
		Where("id = ? AND category_id = ?", subCategoryID, categoryID).
		First(&subCategory).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Sub Category not found",
		})
		return
	}

	// bind request
	var req structs.SubCategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	subCategory.Name = req.Name
	subCategory.Slug = helpers.GenerateSlug(req.Name)

	if err := database.DB.Save(&subCategory).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update sub category",
		})
		return
	}

	// reload relasi
	database.DB.Preload("Category").First(&subCategory, subCategory.ID)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Sub Category updated successfully",
		Data:    subCategory,
	})
}

func DeleteSubCategory(c *gin.Context) {
	categoryID := c.Param("id")
	subCategoryID := c.Param("subCategoryId")

	var subCategory models.SubCategory

	if err := database.DB.
		Where("id = ? AND category_id = ?", subCategoryID, categoryID).
		First(&subCategory).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Sub Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Delete(&subCategory).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete sub category",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Sub Category deleted successfully",
	})
}
