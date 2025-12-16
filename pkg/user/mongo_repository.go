package user

import (
	"MScannot206/shared"
	"MScannot206/shared/entity"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRepository(
	client *mongo.Client,
	dbName string,
) (*UserRepository, error) {

	if client == nil {
		return nil, errors.New("mongo client is null")
	}

	repo := &UserRepository{
		client: client,
	}

	repo.user = client.Database(dbName).Collection(shared.User)

	return repo, nil
}

type UserRepository struct {
	client *mongo.Client

	user *mongo.Collection
}

func (r *UserRepository) FindUserByUids(ctx context.Context, uids []string) ([]*entity.User, []string, error) {
	requestCount := len(uids)
	var users []*entity.User
	newUids := make([]string, 0, requestCount)

	filter := bson.D{
		{Key: "uid", Value: bson.D{{Key: "$in", Value: uids}}},
	}

	cursor, err := r.user.Find(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		if err == mongo.ErrNoDocuments {
			users = []*entity.User{}
		} else {
			return nil, nil, err
		}
	}

	setUsers := make(map[string]struct{}, len(users))
	for _, u := range users {
		setUsers[u.Uid] = struct{}{}
	}

	for _, uid := range uids {
		if _, found := setUsers[uid]; !found {
			newUids = append(newUids, uid)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, nil, err
	}

	return users, newUids, nil
}

func (r *UserRepository) InsertUserByUids(ctx context.Context, uids []string) ([]*entity.User, error) {
	requestCount := len(uids)

	if requestCount == 0 {
		return []*entity.User{}, nil
	}

	newUsers := make([]*entity.User, 0, requestCount)
	writeModels := make([]mongo.WriteModel, 0, requestCount)
	now := time.Now().UTC()

	type userDocument struct {
		*entity.User `bson:",inline"`
		CreatedAt    time.Time
	}

	for _, uid := range uids {
		newUser := entity.NewUser(uid)
		newUsers = append(newUsers, newUser)
		doc := userDocument{
			User:      newUser,
			CreatedAt: now,
		}

		model := mongo.NewInsertOneModel().SetDocument(doc)
		writeModels = append(writeModels, model)
	}

	_, err := r.user.BulkWrite(ctx, writeModels)
	if err != nil {
		return nil, err
	}

	return newUsers, nil
}
