package session

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

// 7주일
const expireAfterSeconds = 3600 * 24 * 7

func NewSessionRepository(
	ctx context.Context,
	client *mongo.Client,
	dbName string,
) (*SessionRepository, error) {

	if client == nil {
		return nil, errors.New("mongo client is null")
	}

	repo := &SessionRepository{
		client:  client,
		session: client.Database(dbName).Collection(shared.UserSession),
	}

	if err := repo.ensureIndexes(ctx); err != nil {
		return nil, err
	}

	return repo, nil
}

type SessionRepository struct {
	client  *mongo.Client
	session *mongo.Collection
}

func (r *SessionRepository) ensureIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "updated_at", Value: 1},
		},
		Options: options.Index().
			SetExpireAfterSeconds(expireAfterSeconds).
			SetName("session_ttl_idx"),
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := r.session.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepository) SaveUserSessions(ctx context.Context, sessions []*entity.UserSession) error {
	if len(sessions) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, len(sessions))

	for i, session := range sessions {
		session.UpdatedAt = time.Now().UTC()

		filter := bson.D{
			{Key: "uid", Value: session.Uid},
		}

		update := bson.D{
			{Key: "$set", Value: session},
		}

		models[i] = mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
	}

	if len(models) > 0 {
		opts := options.BulkWrite().SetOrdered(false)
		_, err := r.session.BulkWrite(ctx, models, opts)
		if err != nil {
			return err
		}
	}

	return nil
}
