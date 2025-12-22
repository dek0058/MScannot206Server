package user

import (
	"MScannot206/shared"
	"MScannot206/shared/entity"
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrUserMongoRepositoryIsNil = errors.New("user mongo repository is null")

func NewUserMongoRepository(
	client *mongo.Client,
	dbName string,
) (*UserMongoRepository, error) {

	if client == nil {
		return nil, errors.New("mongo client is null")
	}

	if dbName == "" {
		return nil, errors.New("database name is empty")
	}

	repo := &UserMongoRepository{
		client: client,
	}

	repo.user = client.Database(dbName).Collection(shared.User)
	repo.characterName = client.Database(dbName).Collection(shared.CharacterName)

	return repo, nil
}

type UserMongoRepository struct {
	client *mongo.Client

	user          *mongo.Collection
	characterName *mongo.Collection
}

func (r *UserMongoRepository) FindUserByUids(ctx context.Context, uids []string) ([]*entity.User, []string, error) {
	requestCount := len(uids)
	var users []*entity.User
	newUids := make([]string, 0, requestCount)

	filter := bson.D{
		{Key: "_id", Value: bson.D{{Key: "$in", Value: uids}}},
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

func (r *UserMongoRepository) InsertUserByUids(ctx context.Context, uids []string) ([]*entity.User, []string, error) {
	requestCount := len(uids)

	if requestCount == 0 {
		return []*entity.User{}, []string{}, nil
	}

	newUsers := make([]*entity.User, requestCount)
	writeModels := make([]mongo.WriteModel, requestCount)
	now := time.Now().UTC()

	type userDocument struct {
		*entity.User `bson:",inline"`
		CreatedAt    time.Time
	}

	for i, uid := range uids {
		newUsers[i] = entity.NewUser(uid)
		doc := userDocument{
			User:      newUsers[i],
			CreatedAt: now,
		}
		writeModels[i] = mongo.NewInsertOneModel().SetDocument(doc)
	}

	_, err := r.user.BulkWrite(ctx, writeModels, options.BulkWrite().SetOrdered(false))
	if err != nil {
		if bulkErr, ok := err.(mongo.BulkWriteException); ok {
			failedUids := make([]string, len(bulkErr.WriteErrors))

			for i, writeErr := range bulkErr.WriteErrors {
				failedUids[i] = newUsers[writeErr.Index].Uid
			}

			// 성공한 유저 필터링
			successfulUsers := make([]*entity.User, 0, requestCount-len(failedUids))
			failedUidSet := make(map[string]struct{}, len(failedUids))
			for _, uid := range failedUids {
				failedUidSet[uid] = struct{}{}
			}

			for _, u := range newUsers {
				if _, found := failedUidSet[u.Uid]; !found {
					successfulUsers = append(successfulUsers, u)
				}
			}

			return successfulUsers, failedUids, nil
		}
		return nil, nil, err
	}

	return newUsers, []string{}, nil
}

func (r *UserMongoRepository) ExistsCharacterNames(ctx context.Context, names []string) (map[string]bool, error) {
	existsMap := make(map[string]bool, len(names))
	for _, name := range names {
		existsMap[name] = false
	}

	filter := bson.D{
		{Key: "_id", Value: bson.D{{Key: "$in", Value: names}}},
	}

	cursor, err := r.characterName.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var charNames []*entity.CharacterName
	if err := cursor.All(ctx, &charNames); err != nil {
		return nil, err
	}

	for _, cn := range charNames {
		existsMap[cn.Name] = true
	}

	return existsMap, nil
}

func (r *UserMongoRepository) GetUsersCharacterSlotCount(ctx context.Context, uids []string) (map[string]int, error) {
	slotCountMap := make(map[string]int, len(uids))

	return slotCountMap, nil
}

func (r *UserMongoRepository) CreateCharacters(ctx context.Context, infos []*UserCreateCharacterInfo) ([]string, []string, error) {
	if len(infos) == 0 {
		return []string{}, []string{}, nil
	}

	successInfos := make(map[string]*UserCreateCharacterInfo, len(infos))
	failedUids := make(map[string]struct{}, 0)

	charNameModels := make([]mongo.WriteModel, len(infos))

	for i, info := range infos {
		successInfos[info.Uid] = info

		charNameDoc := &entity.CharacterName{
			Name:      info.Name,
			CreatedAt: time.Now().UTC(),
		}

		charNameModels[i] = mongo.NewInsertOneModel().SetDocument(charNameDoc)
	}

	_, err := r.characterName.BulkWrite(ctx, charNameModels, options.BulkWrite().SetOrdered(false))
	if err != nil {
		if bulkErr, ok := err.(mongo.BulkWriteException); ok {
			for _, writeErr := range bulkErr.WriteErrors {
				failedUids[infos[writeErr.Index].Uid] = struct{}{}
			}
		} else {
			return nil, nil, err
		}
	}

	for uid := range failedUids {
		delete(successInfos, uid)
	}

	newCharModels := make([]mongo.WriteModel, 0, len(successInfos))
	for _, info := range successInfos {
		char := &entity.Character{
			Slot: info.Slot,
			Name: info.Name,
		}

		update := bson.D{
			{Key: "$push", Value: bson.D{
				{Key: "characters", Value: char},
			}},
		}
		filter := bson.D{
			{Key: "_id", Value: info.Uid},
		}

		newCharModels = append(newCharModels, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(false))
	}

	removeCharNames := make([]string, 0, len(failedUids))
	removeCharNameModels := make([]mongo.WriteModel, 0, len(failedUids))

	var errException error
	if len(newCharModels) > 0 {
		_, err = r.user.BulkWrite(ctx, newCharModels, options.BulkWrite().SetOrdered(false))
		if err != nil {
			if bulkErr, ok := err.(mongo.BulkWriteException); ok {
				errUids := make([]string, 0, len(bulkErr.WriteErrors))
				for _, writeErr := range bulkErr.WriteErrors {
					failedUids[infos[writeErr.Index].Uid] = struct{}{}
					errUids = append(errUids, infos[writeErr.Index].Uid)
				}

				// 캐릭터 생성에 실패하였다면 이전에 기록된 캐릭터 이름도 삭제
				for _, info := range successInfos {
					if _, ok := failedUids[info.Uid]; ok {
						removeCharNames = append(removeCharNames, info.Name)

						opts := mongo.NewDeleteOneModel().SetFilter(bson.D{
							{Key: "_id", Value: info.Name},
						})
						removeCharNameModels = append(removeCharNameModels, opts)
					}
				}

				for _, uid := range errUids {
					delete(successInfos, uid)
				}
			} else {
				// 기타 오류가 발생하였음으로 생성한 캐릭터 이름을 전부 삭제 요청
				log.Warn().Msg("캐릭터 생성 중 오류 발생, 생성된 캐릭터 이름 삭제 시도")
				errException = err
				for _, info := range successInfos {
					removeCharNames = append(removeCharNames, info.Name)

					opts := mongo.NewDeleteOneModel().SetFilter(bson.D{
						{Key: "_id", Value: info.Name},
					})
					removeCharNameModels = append(removeCharNameModels, opts)
				}
			}
		}
	}

	if len(removeCharNameModels) > 0 {
		_, err = r.characterName.BulkWrite(ctx, removeCharNameModels, options.BulkWrite().SetOrdered(false))
		if err != nil {
			if bulkErr, ok := err.(mongo.BulkWriteException); ok {
				log.Warn().Msg("일부 캐릭터 이름 삭제에 실패했습니다")
				for _, writeErr := range bulkErr.WriteErrors {
					log.Warn().Msgf("캐릭터 이름 삭제 실패: %v - %v", removeCharNames[writeErr.Index], writeErr.Message)
				}
			}
			log.Err(err).Msg("캐릭터 이름 삭제 중 오류 발생")
		}
	}

	if errException != nil {
		return nil, nil, errException
	}

	fail := make([]string, 0, len(failedUids))
	for uid := range failedUids {
		fail = append(fail, uid)
	}

	success := make([]string, 0, len(successInfos))
	for uid := range successInfos {
		success = append(success, uid)
	}

	return success, fail, nil
}

func (r *UserMongoRepository) FindCharacters(ctx context.Context, uids []string) (map[string][]*entity.Character, error) {
	charMap := make(map[string][]*entity.Character, len(uids))

	filter := bson.D{
		{Key: "_id", Value: bson.D{{Key: "$in", Value: uids}}},
	}

	opts := options.Find().SetProjection(
		bson.M{
			"_id":        1,
			"characters": 1,
		},
	)

	cursor, err := r.user.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*entity.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	for _, u := range users {
		charMap[u.Uid] = u.Characters
	}

	return charMap, nil
}
