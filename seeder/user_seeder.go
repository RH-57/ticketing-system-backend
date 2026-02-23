package seeders

import (
	"backend-golang-api/database"
	"backend-golang-api/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedSuperAdmin() {

	var count int64
	database.DB.Model(&models.User{}).Count(&count)

	// Kalau sudah ada user, stop
	if count > 0 {
		return
	}

	password, _ := bcrypt.GenerateFromPassword(
		[]byte("admin123"),
		bcrypt.DefaultCost,
	)

	user := models.User{
		Name:     "Super Admin",
		Username: "superadmin",
		Email:    "admin@mail.com",
		Password: string(password),
		Role:     "superadmin",
		IsActive: true,
	}

	database.DB.Create(&user)
}
