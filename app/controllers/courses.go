package controllers

import (
	"TeachAssistApi/app"
	"TeachAssistApi/app/helpers"
	"TeachAssistApi/app/teachassist"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := helpers.AuthenticateUser(c)
		if user == nil {
			return
		}

		metadata, err := helpers.UserToUserMetadata(user)
		if helpers.HandleAppError(err, c) {
			return
		}

		courses, err := teachassist.GetAllCourses(metadata)
		if helpers.HandleAppError(err, c) {
			return
		}

		c.JSON(http.StatusOK, &courses)
	}
}

func GetCourseByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := helpers.AuthenticateUser(c)
		if user == nil {
			return
		}

		metadata, err := helpers.UserToUserMetadata(user)
		if helpers.HandleAppError(err, c) {
			return
		}

		id := c.Param("id")
		if id == "" {
			helpers.HandleAppError(app.CreateError(app.InvalidCourseIdError), c)
			return
		}

		weights, assessments, err := teachassist.GetCourseByID(id, metadata)
		if helpers.HandleAppError(err, c) {
			return
		}

		c.JSON(http.StatusOK, gin.H{"weights": &weights, "assessments": &assessments})
	}
}
