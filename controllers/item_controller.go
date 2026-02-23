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

func ShowItemBySubCategory(c *gin.Context) {
	subCategoryID := c.Param("subCategoryId")

	// cek subcategory
	var subCategory models.SubCategory
	if err := database.DB.First(&subCategory, subCategoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Sub Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var items []models.Item

	if err := database.DB.
		Preload("SubCategory").
		Where("sub_category_id = ?", subCategoryID).
		Order("name ASC").
		Find(&items).Error; err != nil {

		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch items",
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Show all items",
		Data:    items,
	})
}

func CreateItem(c *gin.Context) {
	subCategoryID := c.Param("subCategoryId")

	// cek subcategory
	var subCategory models.SubCategory
	if err := database.DB.First(&subCategory, subCategoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Sub Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var req structs.ItemCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	slug := helpers.GenerateSlug(req.Name)

	var item models.Item

	// cek existing termasuk soft delete
	err := database.DB.Unscoped().
		Where("slug = ? AND sub_category_id = ?", slug, subCategory.ID).
		First(&item).Error

	// RESTORE
	if err == nil && item.DeletedAt.Valid {
		item.Name = req.Name
		item.Slug = slug
		item.DeletedAt = gorm.DeletedAt{}

		if err := database.DB.Unscoped().Save(&item).Error; err != nil {
			c.JSON(http.StatusBadRequest, structs.ErrorResponse{
				Success: false,
				Message: "Failed to restore item",
			})
			return
		}
	} else {
		// CREATE BARU
		item = models.Item{
			Name:          req.Name,
			Slug:          slug,
			SubCategoryID: subCategory.ID,
		}

		if err := database.DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusBadRequest, structs.ErrorResponse{
				Success: false,
				Message: "Failed to create item",
			})
			return
		}
	}

	// reload relasi
	if err := database.DB.
		Preload("SubCategory").
		Preload("SubCategory.Category").
		First(&item, item.ID).Error; err != nil {

		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to load item",
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Item saved successfully",
		Data:    item,
	})
}

func GetItemByID(c *gin.Context) {
	subCategoryID := c.Param("subCategoryId")
	itemID := c.Param("itemId")

	// cek subcategory
	var subCategory models.SubCategory
	if err := database.DB.First(&subCategory, subCategoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Sub Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var item models.Item
	if err := database.DB.
		Preload("SubCategory").
		Preload("SubCategory.Category").
		Where("sub_category_id = ?", subCategory.ID).
		First(&item, itemID).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Item not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Show item detail",
		Data:    item,
	})
}

func UpdateItem(c *gin.Context) {
	subCategoryID := c.Param("subCategoryId")
	itemID := c.Param("itemId")

	// cek subcategory
	var subCategory models.SubCategory
	if err := database.DB.First(&subCategory, subCategoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Sub Category not found",
		})
		return
	}

	// cek item
	var item models.Item
	if err := database.DB.
		Where("sub_category_id = ?", subCategory.ID).
		First(&item, itemID).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Item not found",
		})
		return
	}

	var req structs.ItemUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	slug := helpers.GenerateSlug(req.Name)

	// cek slug duplicate (kecuali dirinya sendiri)
	var existing models.Item
	if err := database.DB.
		Where("slug = ? AND sub_category_id = ? AND id != ?", slug, subCategory.ID, item.ID).
		First(&existing).Error; err == nil {

		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Item with this name already exists",
		})
		return
	}

	item.Name = req.Name
	item.Slug = slug

	if err := database.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update item",
		})
		return
	}

	// reload relasi
	database.DB.
		Preload("SubCategory").
		Preload("SubCategory.Category").
		First(&item, item.ID)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Item updated successfully",
		Data:    item,
	})
}

func DeleteItem(c *gin.Context) {
	subCategoryID := c.Param("subCategoryId")
	itemID := c.Param("itemId")

	// cek subcategory
	var subCategory models.SubCategory
	if err := database.DB.First(&subCategory, subCategoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Sub Category not found",
		})
		return
	}

	// cek item
	var item models.Item
	if err := database.DB.
		Where("sub_category_id = ?", subCategory.ID).
		First(&item, itemID).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Item not found",
		})
		return
	}

	if err := database.DB.Delete(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete item",
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Item deleted successfully",
	})
}
