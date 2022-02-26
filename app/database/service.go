package database

import (
	"TeachAssistApi/app/teachassist"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	DB *mongo.Client
}

func (s *Service) CreateAndUpdateUserIfNecessary(metadata *teachassist.UserMetadata, notifications bool) error {
	user := User{
		Id:            primitive.NewObjectID(),
		Username:      metadata.Username,
		Password:      metadata.Password,
		StudentId:     metadata.StudentId,
		SessionToken:  metadata.SessionToken,
		SessionExpiry: metadata.SessionExpiry,
		Notifications: notifications,
	}

	exists := user.Exists(s.DB)
	if !exists {
		err := user.Create(s.DB)
		return err
	} else {
		err := user.Update(s.DB)
		return err
	}
}
