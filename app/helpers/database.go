package helpers

import (
	"TeachAssistApi/app/database"
	"TeachAssistApi/app/security"
	"TeachAssistApi/app/teachassist"
)

func UserToUserMetadata(user *database.User) (*teachassist.UserMetadata, error) {
	cs, err := security.NewCryptographyService()
	if err != nil {
		return nil, err
	}

	decryptedPassword, err := cs.DecryptFromBase64String(user.Password)
	if err != nil {
		return nil, err
	}
	decryptedToken, err := cs.DecryptFromBase64String(user.SessionToken)
	if err != nil {
		return nil, err
	}

	metadata := &teachassist.UserMetadata{
		Username:      user.Username,
		Password:      decryptedPassword,
		StudentId:     user.StudentId,
		SessionToken:  decryptedToken,
		SessionExpiry: user.SessionExpiry,
	}
	return metadata, nil
}
