package main

import (
	"TeachAssistApi/teachassist"
	"github.com/gin-gonic/gin"
	"net/http"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Group("users/")

	r.POST("login", func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()
		if !hasAuth {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no auth"})
			return
		}

		metadata, err := teachassist.LoginUser(username, password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, metadata)
	})

	return r
}

func main() {
	r := setupRouter()
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
