package controllers

import (
	"backend-golang-api/database"
	"backend-golang-api/helpers"
	"backend-golang-api/models"
	"backend-golang-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowCategories(c *gin.Context) {
	var categories []models.Category

	database.DB.Find(&categories)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List Categories",
		Data:    categories,
	})
}

func CreateCategory(c *gin.Context) {
	var req structs.CategoryCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	category := models.Category{
		Name: req.Name,
		Slug: helpers.GenerateSlug(req.Name),
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create category",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Category created successfully",
		Data:    category,
	})
}

func ShowCategoryById(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Category found",
		Data:    category,
	})
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var req structs.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	category.Name = req.Name
	category.Slug = helpers.GenerateSlug(req.Name)

	database.DB.Save(&category)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Category updated successfully",
		Data:    category,
	})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Delete(&category, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete category",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: false,
		Message: "Category deleted successfully",
	})
}
