package routes

import (
	"TeachAssistApi/app/controllers"
	"github.com/gin-gonic/gin"
)

func AddCoursesRoutes(api *gin.RouterGroup) {
	courses := api.Group("/courses")
	courses.GET("/all", controllers.GetAllCourses())
	courses.GET("/id/:id", controllers.GetCourseByID())
}
