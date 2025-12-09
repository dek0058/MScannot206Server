package repository

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id  primitive.ObjectID `json:"_id" bson:"_id"`
	Uid string             `json:"uid" bson:"uid"`
}

type UserRepository interface {
	FindUserByUID(string) (*User, error)
}
