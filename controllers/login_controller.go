package controllers

import (
	"backend-golang-api/database"
	"backend-golang-api/helpers"
	"backend-golang-api/models"
	"backend-golang-api/structs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var req structs.UserLoginRequest
	var user models.User

	// 1️⃣ Validasi request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// 2️⃣ Cari user berdasarkan email ATAU username
	if err := database.DB.
		Where("email = ? OR username = ?", req.Identifier, req.Identifier).
		First(&user).Error; err != nil {

		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "Invalid email/username or password",
		})
		return
	}

	// 3️⃣ Cek status user (opsional tapi disarankan)
	if user.IsActive != true {
		c.JSON(http.StatusForbidden, structs.ErrorResponse{
			Success: false,
			Message: "User is inactive",
		})
		return
	}

	// 4️⃣ Verifikasi password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {

		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "Invalid email/username or password",
		})
		return
	}

	// 5️⃣ Generate JWT
	accessToken, err := helpers.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	refreshToken, exp, err := helpers.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to generate refresh token",
		})
		return
	}

	session := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: exp,
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
	}

	if err := database.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Faild to store session",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(time.Until(exp).Seconds()),
		"/",
		"localhost",
		false,
		true,
	)

	// 6️⃣ Response sukses
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Login successful",
		Data: structs.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			Avatar:    user.Avatar,
			IsActive:  user.IsActive == true,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
			Token:     &accessToken,
		},
	})
}
