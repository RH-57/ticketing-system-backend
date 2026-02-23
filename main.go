package main

import (
	"backend-golang-api/config"
	"backend-golang-api/database"
	"backend-golang-api/routes"
)

func main() {

	config.LoadEnv()
	database.InitDB()
	r := routes.SetupRouter()

	r.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
