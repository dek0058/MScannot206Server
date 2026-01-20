package channel

import (
	"MScannot206/shared"
	"MScannot206/shared/entity"
	"context"
	"errors"
	"time"

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

type channelDoc struct {
	ID       string    `bson:"_id"`
	CreateAt time.Time `bson:"create_at"`
}

func (r *ChannelMongoRepository) AddChannel(ctx context.Context, id string) error {

	doc := &channelDoc{
		ID:       id,
		CreateAt: time.Now(),
	}

	_, err := r.channel.InsertOne(ctx, doc)

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
	opts := options.Find().SetSort(
		bson.D{{Key: "create_at", Value: -1}},
	)

	cursor, err := r.channel.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	docs := make([]*channelDoc, 0)
	for cursor.Next(ctx) {
		var doc channelDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		docs = append(docs, &doc)
	}

	channels := make([]*entity.Channel, 0, len(docs))
	for i, doc := range docs {
		channels = append(channels, &entity.Channel{
			Id:    doc.ID,
			Index: i + 1,
		})
	}

	return channels, nil
}
