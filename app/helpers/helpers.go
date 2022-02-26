package helpers

import (
	"TeachAssistApi/app"
	"TeachAssistApi/app/security"
	"github.com/gin-gonic/gin"
	"strings"
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

func ExtractBearerToken(c *gin.Context) (*security.ParsedToken, error) {
	authorization := c.Request.Header["Authorization"]
	if len(authorization) == 0 {
		return nil, app.CreateError(app.AuthError)
	}

	bearer := authorization[0]
	if bearer == "" || !strings.Contains(bearer, "Bearer") {
		return nil, app.CreateError(app.AuthError)
	}

	separated := strings.Split(bearer, " ")
	if len(separated) < 2 {
		return nil, app.CreateError(app.AuthError)
	}

	token := separated[1]

	parsed := security.VerifyJWT(token)

	return &parsed, nil
}
