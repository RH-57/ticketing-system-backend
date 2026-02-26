package controllers

import (
	"backend-golang-api/database"
	"backend-golang-api/helpers"
	"backend-golang-api/models"
	"backend-golang-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllDepartment(c *gin.Context) {
	var departments []models.Department

	err := database.DB.
		Preload("Division").
		Find(&departments).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed get departments",
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Get All Departments",
		Data:    departments,
	})
}

func ShowDepartmentByDivision(c *gin.Context) {
	branchId := c.Param("id")
	divisionId := c.Param("divisionId")

	var division models.Division

	// validasi division berdasarkan branch
	if err := database.DB.
		Where("id = ? AND branch_id = ?", divisionId, branchId).
		First(&division).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Division not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var departments []models.Department
	if err := database.DB.
		Where("division_id = ?", divisionId).
		Order("created_at ASC").
		Find(&departments).Error; err != nil {

		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch Departments",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Show all Departments",
		Data:    departments,
	})
}

func CreateDepartment(c *gin.Context) {
	branchId := c.Param("id")
	divisionId := c.Param("divisionId")

	var division models.Division
	if err := database.DB.
		Where("id = ? AND branch_id = ?", divisionId, branchId).
		First(&division).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Division not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var input structs.DepartmentCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Invalid request",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	department := models.Department{
		Name:       input.Name,
		DivisionID: division.ID,
	}

	if err := database.DB.Create(&department).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create Department",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Department created successfully",
		Data:    department,
	})
}

func UpdateDepartment(c *gin.Context) {
	branchId := c.Param("id")
	divisionId := c.Param("divisionId")
	departmentId := c.Param("departmentId")

	var department models.Department
	if err := database.DB.
		Joins("JOIN divisions ON divisions.id = departments.division_id").
		Where("departments.id = ? AND divisions.id = ? AND divisions.branch_id = ?",
			departmentId, divisionId, branchId).
		First(&department).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Department not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var input structs.DepartmentUpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Invalid request",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	department.Name = input.Name

	if err := database.DB.Save(&department).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update Department",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Department updated successfully",
		Data:    department,
	})
}

func DeleteDepartment(c *gin.Context) {
	branchId := c.Param("id")
	divisionId := c.Param("divisionId")
	departmentId := c.Param("departmentId")

	var department models.Department
	if err := database.DB.
		Joins("JOIN divisions ON divisions.id = departments.division_id").
		Where("departments.id = ? AND divisions.id = ? AND divisions.branch_id = ?",
			departmentId, divisionId, branchId).
		First(&department).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Department not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Delete(&department).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete Department",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Department deleted successfully",
	})
}
