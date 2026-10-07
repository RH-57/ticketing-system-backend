package main

import (
	"backend-golang-api/config"
	"backend-golang-api/database"
	"backend-golang-api/routes"
	seeders "backend-golang-api/seeder"
)

func main() {

	config.LoadEnv()
	database.InitDB()

	seeders.SeedSuperAdmin()
	r := routes.SetupRouter()

	r.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
