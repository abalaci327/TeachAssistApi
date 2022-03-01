package responses

type LoginUserResponse struct {
	Token string `json:"token" example:"jwt"`
}

type DeleteUserResponse struct {
	Success bool `json:"success" example:"true"`
}
