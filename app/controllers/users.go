package controllers

import (
	"TeachAssistApi/app"
	"TeachAssistApi/app/controllers/responses"
	"TeachAssistApi/app/database"
	"TeachAssistApi/app/security"
	"TeachAssistApi/app/teachassist"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()
		if !hasAuth {
			err := app.CreateError(app.AuthError)
			c.JSON(http.StatusBadRequest, err.ErrorResponse())
			return
		}
		q := c.Request.URL.Query().Get("notifications")
		notifications := q == "true"

		metadata, err := teachassist.LoginUser(username, password)
		if err != nil {
			if e, ok := (err).(app.Error); ok {
				c.JSON(e.StatusCode, e.ErrorResponse())
			}
			return
		}

		user := database.User{
			Id:            primitive.NewObjectID(),
			Username:      metadata.Username,
			Password:      metadata.Password,
			StudentId:     metadata.StudentId,
			SessionToken:  metadata.SessionToken,
			SessionExpiry: metadata.SessionExpiry,
			Notifications: notifications,
		}

		err = user.Create(database.DB)
		if err != nil {
			if e, ok := (err).(app.Error); ok {
				c.JSON(e.StatusCode, e.ErrorResponse())
			}
			return
		}

		jwt, err := security.CreateJWT(metadata.Username, metadata.StudentId, notifications)
		if err != nil {
			if e, ok := (err).(app.Error); ok {
				c.JSON(e.StatusCode, e.ErrorResponse())
			}
			return
		}

		c.JSON(http.StatusCreated, responses.LoginUserResponse{Token: jwt})
	}
}
