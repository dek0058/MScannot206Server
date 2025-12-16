package serverinfo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "service_configs"

func NewServerInfoRepository(client *mongo.Client, dbName string) (*ServerInfoRepository, error) {
	if client == nil {
		return nil, errors.New("mongo client is null")
	}

	return &ServerInfoRepository{
		client:     client,
		collection: client.Database(dbName).Collection(collectionName),
	}, nil
}

type ServerInfoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (r *ServerInfoRepository) GetConfig(ctx context.Context, name string) (*ServerInfo, error) {
	filter := bson.M{"_id": name}

	var info ServerInfo
	err := r.collection.FindOne(ctx, filter).Decode(&info)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &info, nil
}

func (r *ServerInfoRepository) SaveConfig(ctx context.Context, info *ServerInfo) error {
	info.UpdatedAt = time.Now().UTC()

	filter := bson.M{"_id": info.Name}
	update := bson.M{"$set": info}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *ServerInfoRepository) UpdateStatus(ctx context.Context, name string, status ServerStatus) error {
	filter := bson.M{"_id": name}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now().UTC(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
