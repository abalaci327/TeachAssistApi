package main

import (
	"TeachAssistApi/app/database"
	"TeachAssistApi/app/helpers"
	"TeachAssistApi/app/routes"
	_ "TeachAssistApi/docs"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	api := r.Group("/api/v1")
	routes.AddUserRoutes(api)
	routes.AddCoursesRoutes(api)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// @title                       TeachAssist API
// @version                     0.1.0
// @description                 The fast easy and simple way to access all of your YRDSB course marks.
// @host                        localhost:8080
// @BasePath                    /api/v1
// @securityDefinitions.basic   BasicAuth
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
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
