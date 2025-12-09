package mongodb

import (
	"MScannot206/shared/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRepository(
	client *mongo.Client,
) (*UserRepository, error) {

	// if client == nil {
	// 	return nil, errors.New("mongo client is null")
	// }

	return &UserRepository{
		client: client,
	}, nil
}

type UserRepository struct {
	client *mongo.Client

	collection *mongo.Collection
}

func (r *UserRepository) FindUserByUID(uid string) (*repository.User, error) {

	return nil, nil
}
