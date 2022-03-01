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

// LoginUser godoc
// @Summary      Login User
// @Description  Login a TeachAssist YRDSB user and optionally configure notifications.
// @Tags         Users
// @Produce      json
// @Param        notifications  query     bool  true  "Enable Notifications"
// @Success      200            {object}  responses.LoginUserResponse
// @Failure      400            {object}  app.ErrorResponse
// @Failure      401            {object}  app.ErrorResponse
// @Failure      500            {object}  app.ErrorResponse
// @Failure      502            {object}  app.ErrorResponse
// @Router       /users/login [post]
// @Security     BasicAuth
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

// RenewUserSession godoc
// @Summary      Renew User Session
// @Description  Get new JWT Token and resubscribe to notifications.
// @Tags         Users
// @Produce      json
// @Param        notifications  path      bool  true  "Enable Notifications"
// @Success      200            {object}  responses.LoginUserResponse
// @Failure      400            {object}  app.ErrorResponse
// @Failure      401            {object}  app.ErrorResponse
// @Failure      500            {object}  app.ErrorResponse
// @Failure      502            {object}  app.ErrorResponse
// @Router       /users/renew_session [get]
// @Security     BasicAuth
func RenewUserSession() gin.HandlerFunc {
	return LoginUser()
}

// RemoveUser godoc
// @Summary      Remove User
// @Description  Remove all user data from the database and revoke JWT Token.
// @Tags         Users
// @Produce      json
// @Success      200  {object}  responses.DeleteUserResponse
// @Failure      400  {object}  app.ErrorResponse
// @Failure      401  {object}  app.ErrorResponse
// @Failure      500  {object}  app.ErrorResponse
// @Failure      502  {object}  app.ErrorResponse
// @Router       /users/remove [delete]
// @Security     BearerAuth
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

		c.JSON(200, responses.DeleteUserResponse{Success: true})
	}
}
