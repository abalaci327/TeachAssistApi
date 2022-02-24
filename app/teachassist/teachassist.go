package teachassist

func LoginUser(username, password string) (UserMetadata, error) {
	return loginUser(username, password)
}
