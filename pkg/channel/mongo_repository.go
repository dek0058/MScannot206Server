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
	client *mongo.Client,
	dbName string,
) (*ChannelMongoRepository, error) {
	if client == nil {
		return nil, ErrChannelRepositoryIsNil
	}

	repo := &ChannelMongoRepository{
		client: client,

		channel: client.Database(dbName).Collection(shared.Channel),
		counter: client.Database(dbName).Collection(shared.Counter),
	}

	return repo, nil
}

type ChannelMongoRepository struct {
	client *mongo.Client

	channel *mongo.Collection
	counter *mongo.Collection
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

func (r *ChannelMongoRepository) CreateLease(ctx context.Context, channel entity.Channel) error {
	_, err := r.channel.InsertOne(ctx, channel)
	return err
}

func (r *ChannelMongoRepository) RenewLease(ctx context.Context, leaseID string, newExpiry time.Time) (*entity.Channel, error) {
	filter := bson.M{
		"lease_id": leaseID,
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

func (r *ChannelMongoRepository) FindExpiredLeases(ctx context.Context, now time.Time) ([]*entity.Channel, error) {

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

func (r *ChannelMongoRepository) DeleteLeases(ctx context.Context, leaseIDs []string) (int64, error) {
	filter := bson.M{
		"lease_id": bson.M{
			"$in": leaseIDs,
		},
	}

	result, err := r.channel.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (r *ChannelMongoRepository) GetAllActiveLeases(ctx context.Context) ([]*entity.Channel, error) {
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
