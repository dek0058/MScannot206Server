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
var SequenceName = "channel"

func NewChannelMongoRepository(
	ctx context.Context,
	client *mongo.Client,
	dbName string,
) (*ChannelMongoRepository, error) {
	if client == nil {
		return nil, ErrChannelRepositoryIsNil
	}

	repo := &ChannelMongoRepository{
		client: client,

		counter: client.Database(dbName).Collection(shared.Counter),

		channel:        client.Database(dbName).Collection(shared.Channel),
		channelRecycle: client.Database(dbName).Collection(shared.ChannelRecycle),
	}

	if err := repo.ensureIndexes(ctx); err != nil {
		return nil, err
	}

	return repo, nil
}

type ChannelMongoRepository struct {
	client *mongo.Client

	counter *mongo.Collection

	channel        *mongo.Collection
	channelRecycle *mongo.Collection
}

func (r *ChannelMongoRepository) ensureIndexes(ctx context.Context) error {
	index := mongo.IndexModel{
		Keys: bson.D{
			{Key: "expires_at", Value: 1},
		},
	}

	_, err := r.channel.Indexes().CreateOne(ctx, index)
	return err
}

func (r *ChannelMongoRepository) GetNextSequence(ctx context.Context) (int, error) {
	filter := bson.M{
		"_id": SequenceName,
	}

	update := bson.M{
		"$inc": bson.M{
			"seq": 1,
		},
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var entity entity.Counter
	err := r.counter.FindOneAndUpdate(ctx, filter, update, opts).Decode(&entity)

	return entity.Seq, err
}

func (r *ChannelMongoRepository) CreateChannel(ctx context.Context, channel entity.Channel) error {
	_, err := r.channel.InsertOne(ctx, channel)
	return err
}

func (r *ChannelMongoRepository) RenewChannel(ctx context.Context, channelId string, newExpiry time.Time) (*entity.Channel, error) {
	filter := bson.M{
		"_id": channelId,
	}

	update := bson.M{
		"$set": bson.M{
			"expires_at": newExpiry,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var entity entity.Channel
	err := r.channel.FindOneAndUpdate(ctx, filter, update, opts).Decode(&entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Lease not found
		}
		return nil, err
	}

	return &entity, nil
}

func (r *ChannelMongoRepository) FindChannelByID(ctx context.Context, channelId string) (*entity.Channel, error) {
	filter := bson.M{
		"_id": channelId,
	}

	var entity entity.Channel
	err := r.channel.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Channel not found
		}
		return nil, err
	}

	return &entity, nil
}

func (r *ChannelMongoRepository) FindExpiredChannels(ctx context.Context, now time.Time) ([]*entity.Channel, error) {

	filter := bson.M{
		"expires_at": bson.M{
			"$lte": now,
		},
	}

	cursor, err := r.channel.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var expired []*entity.Channel
	if err := cursor.All(ctx, &expired); err != nil {
		return nil, err
	}

	return expired, nil
}

func (r *ChannelMongoRepository) DeleteChannels(ctx context.Context, channelIDs []string) (int64, error) {
	filter := bson.M{
		"_id": bson.M{
			"$in": channelIDs,
		},
	}

	result, err := r.channel.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (r *ChannelMongoRepository) GetAllActiveChannels(ctx context.Context) ([]*entity.Channel, error) {
	cursor, err := r.channel.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var activeChannels []*entity.Channel
	if err := cursor.All(ctx, &activeChannels); err != nil {
		return nil, err
	}

	return activeChannels, nil
}

func (r *ChannelMongoRepository) PopRecyclableIndex(ctx context.Context) (int, error) {
	var entity entity.ChannelRecycler
	opts := options.FindOneAndDelete().SetSort(bson.M{"index": 1})
	err := r.channelRecycle.FindOneAndDelete(ctx, bson.M{}, opts).Decode(&entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil // No recyclable index available
		}
		return 0, err
	}
	return entity.Index, nil
}

func (r *ChannelMongoRepository) PushRecyclableIndex(ctx context.Context, index int) error {
	recycler := entity.ChannelRecycler{
		Index: index,
	}
	_, err := r.channelRecycle.InsertOne(ctx, recycler)
	return err
}
