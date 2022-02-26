package database

import (
	"TeachAssistApi/app"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type User struct {
	Id            primitive.ObjectID
	Username      string    `bson:"username,omitempty" validate:"required"`
	Password      string    `bson:"password,omitempty" validate:"required"`
	StudentId     string    `bson:"student_id,omitempty" validate:"required"`
	SessionToken  string    `bson:"session_token,omitempty" validate:"required"`
	SessionExpiry time.Time `bson:"session_expiry" validate:"required"`
	Notifications bool      `bson:"notifications" validate:"required"`
}

func (u *User) Create(db *mongo.Client) error {
	users := db.Database("teachassist").Collection("users")
	_, err := users.InsertOne(context.TODO(), u)
	if err != nil {
		return app.CreateError(app.DatabaseError)
	}
	return nil
}

func (u *User) Read(db *mongo.Client) error {
	users := db.Database("teachassist").Collection("users")
	found := users.FindOne(context.TODO(), bson.M{"username": u.Username})
	err := found.Decode(&u)
	if err != nil {
		return app.CreateError(app.DatabaseError)
	}
	return nil
}

func (u *User) Exists(db *mongo.Client) bool {
	users := db.Database("teachassist").Collection("users")
	found := users.FindOne(context.TODO(), bson.M{"username": u.Username})
	var user *User
	err := found.Decode(&user)
	if err != nil {
		return false
	}
	return true
}

func (u *User) Update(db *mongo.Client) error {
	users := db.Database("teachassist").Collection("users")
	update := bson.M{"username": u.Username, "password": u.Password, "student_id": u.StudentId, "session_token": u.SessionToken, "session_expiry": u.SessionExpiry, "notifications": u.Notifications}
	result := users.FindOneAndUpdate(context.TODO(), bson.M{"username": u.Username}, bson.M{"$set": update})
	if result.Err() != nil {
		return app.CreateError(app.DatabaseError)
	}
	return nil
}

func (u *User) Delete(db *mongo.Client) error {
	users := db.Database("teachassist").Collection("users")
	_, err := users.DeleteOne(context.TODO(), bson.M{"username": u.Username})
	if err != nil {
		return app.CreateError(app.DatabaseError)
	}
	return nil
}
