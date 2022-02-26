package helpers

import (
	"TeachAssistApi/app"
	"github.com/gin-gonic/gin"
)

func HandleAppError(err error, c *gin.Context) bool {
	if err != nil {
		if e, ok := (err).(app.Error); ok {
			c.JSON(e.StatusCode, e.ErrorResponse())
		} else {
			err = app.CreateError(app.UnknownError)
			c.JSON(e.StatusCode, e.ErrorResponse())
		}
		return true
	}
	return false
}
