package channel

import (
	"MScannot206/shared"
	"MScannot206/shared/entity"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrChannelRepositoryIsNil = errors.New("channel repository is nil")

func NewChannelMongoRepository(
	client *mongo.Client,
	dbName string,
) (*ChannelMongoRepository, error) {
	if client == nil {
		return nil, ErrChannelRepositoryIsNil
	}

	repo := &ChannelMongoRepository{
		client: client,

		channel: client.Database(dbName).Collection(shared.Channel),
	}

	return repo, nil
}

type ChannelMongoRepository struct {
	client *mongo.Client

	channel *mongo.Collection
}

func (r *ChannelMongoRepository) AddChannel(ctx context.Context, id string) error {
	filter := bson.M{
		"_id": id,
	}

	update := bson.M{
		"$inc": bson.M{"index": 1},
		"$setOnInsert": bson.M{
			"_id": id,
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := r.channel.UpdateOne(ctx, filter, update, opts)

	return err
}

func (r *ChannelMongoRepository) RemoveChannel(ctx context.Context, id string) error {
	filter := bson.M{
		"_id": id,
	}

	_, err := r.channel.DeleteOne(ctx, filter)

	return err
}

func (r *ChannelMongoRepository) GetChannels(ctx context.Context) ([]*entity.Channel, error) {
	cursor, err := r.channel.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	channels := make([]*entity.Channel, 0)
	for cursor.Next(ctx) {
		var channel entity.Channel
		if err := cursor.Decode(&channel); err != nil {
			return nil, err
		}
		channels = append(channels, &channel)
	}

	return channels, nil
}
