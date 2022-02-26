package controllers

import (
	"TeachAssistApi/app"
	"TeachAssistApi/app/controllers/responses"
	"TeachAssistApi/app/database"
	"TeachAssistApi/app/helpers"
	"TeachAssistApi/app/security"
	"TeachAssistApi/app/teachassist"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()
		if !hasAuth {
			err := app.CreateError(app.AuthError)
			if helpers.HandleAppError(err, c) {
				return
			}
		}
		q := c.Request.URL.Query().Get("notifications")
		notifications := q == "true"

		metadata, err := teachassist.LoginUser(username, password)
		if helpers.HandleAppError(err, c) {
			return
		}

		s := database.Service{DB: database.DB}
		err = s.CreateAndUpdateUserIfNecessary(&metadata, notifications)
		if helpers.HandleAppError(err, c) {
			return
		}

		jwt, err := security.CreateJWT(metadata.Username, metadata.StudentId, notifications)
		if helpers.HandleAppError(err, c) {
			return
		}

		c.JSON(http.StatusCreated, responses.LoginUserResponse{Token: jwt})
	}
}

func RenewUserSession() gin.HandlerFunc {
	return LoginUser()
}

func RemoveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := helpers.ExtractBearerToken(c)
		if helpers.HandleAppError(err, c) {
			return
		}
		if !token.Valid {
			helpers.HandleAppError(app.CreateError(app.AuthError), c)
			return
		}

		user := database.User{Username: token.Username}

		s := database.Service{DB: database.DB}
		err = s.DeleteUser(&user)
		if helpers.HandleAppError(err, c) {
			return
		}

		c.JSON(200, responses.DeleteUserResponse)
	}
}
