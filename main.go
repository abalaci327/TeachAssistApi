package main

import (
	"TeachAssistApi/app"
	"TeachAssistApi/app/database"
	"TeachAssistApi/app/helpers"
	"TeachAssistApi/app/teachassist"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Group("users/")

	r.POST("login", func(c *gin.Context) {
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

		c.JSON(http.StatusOK, metadata)
	})

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
