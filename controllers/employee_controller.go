package controllers

import (
	"backend-golang-api/database"
	"backend-golang-api/helpers"
	"backend-golang-api/models"
	"backend-golang-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowEmployees(c *gin.Context) {
	var employees []models.Employee

	database.DB.
		Preload("Branch").
		Preload("Division").
		Preload("Department").
		Find(&employees)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List all employees",
		Data:    employees,
	})
}

func CreateEmployee(c *gin.Context) {
	var req structs.EmployeeCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var department models.Department
	if err := database.DB.First(&department, req.DepartmentID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Department not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var division models.Division
	if err := database.DB.First(&division, department.DivisionID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Division not found",
		})
		return
	}

	employee := models.Employee{
		Name:         req.Name,
		DepartmentID: department.ID,
		DivisionID:   division.ID,
		BranchID:     division.BranchID,
	}

	if err := database.DB.Create(&employee).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create employee",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// preload supaya response lengkap
	database.DB.
		Preload("Branch").
		Preload("Division").
		Preload("Department").
		First(&employee, employee.ID)

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Employee created successfully",
		Data:    employee,
	})
}

func ShowEmployeeById(c *gin.Context) {
	id := c.Param("id")

	var employee []models.Employee

	if err := database.DB.
		Preload("Branch").
		Preload("Division").
		Preload("Department").
		First(&employee, id).Error; err != nil {

		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Employee not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Employee found",
		Data:    employee,
	})
}

func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")

	var req structs.EmployeeUpdateRequest
	var employee models.Employee

	if err := database.DB.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Employee not found",
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

	// ambil department baru
	var department models.Department
	if err := database.DB.First(&department, req.DepartmentID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Department not found",
		})
		return
	}

	// ambil division dari department
	var division models.Division
	if err := database.DB.First(&division, department.DivisionID).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Division not found",
		})
		return
	}

	// update data
	employee.Name = req.Name
	employee.DepartmentID = department.ID
	employee.DivisionID = division.ID
	employee.BranchID = division.BranchID

	if err := database.DB.Save(&employee).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update employee",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	database.DB.
		Preload("Branch").
		Preload("Division").
		Preload("Department").
		First(&employee, employee.ID)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Employee updated successfully",
		Data:    employee,
	})
}

func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	var employee models.Employee

	if err := database.DB.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Employee not found",
		})
		return
	}

	if err := database.DB.Delete(&employee).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete employee",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Employee deleted successfully",
	})
}
