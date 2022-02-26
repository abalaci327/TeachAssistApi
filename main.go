package main

import (
	"TeachAssistApi/app/database"
	"TeachAssistApi/app/helpers"
	"TeachAssistApi/app/routes"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	api := r.Group("/api/v1")
	routes.AddUserRoutes(api)

	return r
}

func main() {
	helpers.LoadEnvironment()

	database.ConnectDatabase()
	defer database.DisconnectDatabase(database.DB)

	r := setupRouter()
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
