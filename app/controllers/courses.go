package controllers

import (
	"TeachAssistApi/app/helpers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := helpers.AuthenticateUser(c)
		if user == nil {
			return
		}
		c.JSON(http.StatusOK, gin.H{"all": []string{}})
	}
}

func GetCourseByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := helpers.AuthenticateUser(c)
		if user == nil {
			return
		}
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{"id": id})
	}
}
