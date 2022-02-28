package helpers

import (
	"TeachAssistApi/app/database"
	"TeachAssistApi/app/teachassist"
)

func UserToUserMetadata(user *database.User) *teachassist.UserMetadata {
	metadata := &teachassist.UserMetadata{
		Username:      user.Username,
		Password:      user.Password,
		StudentId:     user.StudentId,
		SessionToken:  user.SessionToken,
		SessionExpiry: user.SessionExpiry,
	}
	return metadata
}
