package database

import (
	"backend-golang-api/config"
	"backend-golang-api/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	// Load konfigurasi database dari .env
	dbUser := config.GetEnv("DB_USER", "postgres")
	dbPass := config.GetEnv("DB_PASS", "")
	dbHost := config.GetEnv("DB_HOST", "localhost")
	dbPort := config.GetEnv("DB_PORT", "5432")
	dbName := config.GetEnv("DB_NAME", "db_ticketing")
	dbSSL := config.GetEnv("DB_SSL", "disable")

	// ✅ Format DSN PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		dbHost,
		dbUser,
		dbPass,
		dbName,
		dbPort,
		dbSSL,
	)

	// Koneksi ke database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("PostgreSQL connected successfully!")

	//err = DB.Migrator().DropTable()
	//if err != nil {
	//log.Fatal("Failed dropping refresh_tokens:", err)
	//}

	// Auto Migrate Models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Branch{},
		&models.Division{},
		&models.Department{},
		&models.Employee{},
		&models.Category{},
		&models.SubCategory{},
		&models.Item{},
		&models.Ticket{},
		&models.Comment{},
		&models.RefreshToken{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database migrated successfully!")
}
