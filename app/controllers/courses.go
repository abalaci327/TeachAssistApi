package controllers

import (
	"TeachAssistApi/app"
	"TeachAssistApi/app/controllers/responses"
	"TeachAssistApi/app/helpers"
	"TeachAssistApi/app/teachassist"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAllCourses godoc
// @Summary      Get All Courses
// @Description  Get metadata about all courses a user is currently enrolled in.
// @Tags         Courses
// @Produce      json
// @Success      200  {object}  responses.AllCoursesResponse
// @Failure      400  {object}  app.ErrorResponse
// @Failure      401  {object}  app.ErrorResponse
// @Failure      500  {object}  app.ErrorResponse
// @Failure      502  {object}  app.ErrorResponse
// @Router       /courses/all [get]
// @Security     BearerAuth
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

		c.JSON(http.StatusOK, responses.AllCoursesResponse{Metadata: *courses})
	}
}

// GetCourseByID godoc
// @Summary      Get Course By ID
// @Description  Gets the course with the provided ID including mark weightings and assignments.
// @Tags         Courses
// @Produce      json
// @Success      200  {object}  responses.CourseIDResponse
// @Failure      400  {object}  app.ErrorResponse
// @Failure      401  {object}  app.ErrorResponse
// @Failure      500  {object}  app.ErrorResponse
// @Failure      502  {object}  app.ErrorResponse
// @Router       /courses/id/{id} [get]
// @Security     BearerAuth
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

		c.JSON(http.StatusOK, responses.CourseIDResponse{Assessments: *assessments, Weights: *weights})
	}
}
