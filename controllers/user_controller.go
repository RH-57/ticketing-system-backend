package controllers

import (
	"backend-golang-api/database"
	"backend-golang-api/helpers"
	"backend-golang-api/models"
	"backend-golang-api/structs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ShowUsers(c *gin.Context) {
	var users []models.User

	database.DB.Find(&users)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List All Users",
		Data:    users,
	})
}

func CreateUser(c *gin.Context) {

	var req structs.UserCreateRequest

	// 1️⃣ Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// 2️⃣ Create user (tanpa avatar)
	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
		Role:     req.Role,
		IsActive: true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
		})
		return
	}

	// 3️⃣ Response
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "User created successfully",
		Data: structs.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			Avatar:    user.Avatar, // kosong ("" / null)
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	})
}

func FindUserById(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User Not Found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User Found",
		Data: structs.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req structs.UserUpdateRequest
	var user models.User

	// 1️⃣ Cari user
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// 2️⃣ Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// 3️⃣ Update field
	user.Name = req.Name
	user.Username = req.Username
	user.Email = req.Email
	user.Role = req.Role

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// 4️⃣ Update password (opsional)
	if req.Password != "" {
		user.Password = helpers.HashPassword(req.Password)
	}

	// 5️⃣ Save
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update user",
		})
		return
	}

	// 6️⃣ Response
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User updated successfully",
		Data: structs.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	// 1️⃣ Cari user
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// 2️⃣ Delete (soft delete)
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// 3️⃣ Response
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}

func ChangePassword(c *gin.Context) {

	var user models.User
	var req structs.ChangePasswordRequest

	// 🔥 1️⃣ Ambil user_id dari JWT
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// 2️⃣ Cari user
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	// 3️⃣ Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// 4️⃣ Check current password
	if !helpers.CheckPasswordHash(req.CurrentPassword, user.Password) {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Current password is incorrect",
		})
		return
	}

	// 🔥 Optional security: jangan boleh sama dengan password lama
	if helpers.CheckPasswordHash(req.NewPassword, user.Password) {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "New password cannot be same as old password",
		})
		return
	}

	// 5️⃣ Hash password baru
	user.Password = helpers.HashPassword(req.NewPassword)

	// 6️⃣ Save
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to change password",
		})
		return
	}

	// 7️⃣ Response
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}
