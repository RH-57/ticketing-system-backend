package controllers

import (
	"backend-golang-api/database"
	"backend-golang-api/helpers"
	"backend-golang-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RefreshToken(c *gin.Context) {

	// 1️⃣ Ambil dari cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "refresh token missing"})
		return
	}

	// 2️⃣ Parse JWT refresh token
	claims, err := helpers.ParseToken(refreshToken)
	if err != nil || claims.Type != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid refresh token"})
		return
	}

	// 3️⃣ Cek DB session
	var session models.RefreshToken
	err = database.DB.
		Where("token = ? AND revoked = false", refreshToken).
		First(&session).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "session not found"})
		return
	}

	// 4️⃣ Expired?
	if session.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "refresh expired"})
		return
	}

	// 5️⃣ ROTATE SESSION (security best practice)
	session.Revoked = true
	database.DB.Save(&session)

	// 6️⃣ Generate token baru
	accessToken, _ := helpers.GenerateToken(claims.UserID, claims.Email, claims.Role)
	newRefresh, exp, _ := helpers.GenerateRefreshToken(claims.UserID, claims.Email, claims.Role)

	// 7️⃣ Simpan session baru
	database.DB.Create(&models.RefreshToken{
		UserID:    claims.UserID,
		Token:     newRefresh,
		ExpiresAt: exp,
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
	})

	// 8️⃣ Set cookie baru (IMPORTANT)
	c.SetCookie(
		"refresh_token",
		newRefresh,
		int(time.Until(exp).Seconds()),
		"/",
		"",
		false, // dev http
		true,  // httpOnly
	)

	// 9️⃣ Kirim access token saja
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}
