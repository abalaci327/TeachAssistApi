package responses

import "github.com/gin-gonic/gin"

type LoginUserResponse struct {
	Token string `json:"token"`
}

var DeleteUserResponse = gin.H{"success": true}
