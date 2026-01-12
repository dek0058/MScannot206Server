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
	ctx context.Context,
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

		user:          client.Database(dbName).Collection(shared.User),
		characterName: client.Database(dbName).Collection(shared.CharacterName),
	}

	if err := repo.ensureIndexes(ctx); err != nil {
		return nil, err
	}

	return repo, nil
}

type UserMongoRepository struct {
	client *mongo.Client

	user          *mongo.Collection
	characterName *mongo.Collection
}

func (r *UserMongoRepository) ensureIndexes(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	slotIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "characters.slot", Value: 1},
		},
		Options: options.Index().
			SetName("user_character_slot_idx"),
	}

	// 캐릭터 인덱스
	characterNameIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "characters.name", Value: 1},
		},
		Options: options.Index().
			SetName("user_character_name_idx"),
	}

	// 기존 인덱스 삭제 - https://github.com/dek0058/MScannot206Server/issues/5
	if _, err := r.user.Indexes().DropOne(ctx, "user_character_name_idx"); err != nil {
		var cmdErr mongo.CommandError
		if errors.As(err, &cmdErr) {
			if cmdErr.Code != 27 {
				return err
			}
		} else {
			return err
		}
	}

	_, err := r.user.Indexes().CreateMany(ctx, []mongo.IndexModel{slotIndex, characterNameIndex})
	if err != nil {
		return err
	}

	return nil
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
			log.Err(err).Msg("일부 유저 생성에 실패했습니다")
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

func (r *UserMongoRepository) CreateCharacters(ctx context.Context, infos []*UserCreateCharacter) (map[string]*entity.Character, map[string]string, error) {
	if len(infos) == 0 {
		return map[string]*entity.Character{}, map[string]string{}, nil
	}

	failureUids := make(map[string]string, len(infos))
	userRequests := make(map[string]*UserCreateCharacter, len(infos))
	charNameModels := make([]mongo.WriteModel, len(infos))

	for i, info := range infos {
		userRequests[info.Uid] = info

		doc := &entity.CharacterName{
			Name:      info.Name,
			CreatedAt: time.Now().UTC(),
		}
		charNameModels[i] = mongo.NewInsertOneModel().SetDocument(doc)
	}

	_, err := r.characterName.BulkWrite(ctx, charNameModels, options.BulkWrite().SetOrdered(false))
	if err != nil {
		if bulkErr, ok := err.(mongo.BulkWriteException); ok {
			for _, writeErr := range bulkErr.WriteErrors {
				uid := infos[writeErr.Index].Uid
				if mongo.IsDuplicateKeyError(writeErr) {
					failureUids[uid] = USER_CHARACTER_NAME_ALREADY_EXISTS_ERROR
				}
				delete(userRequests, infos[writeErr.Index].Uid)
			}
		} else {
			return nil, nil, err
		}
	}

	createInfos := make([]*UserCreateCharacter, 0, len(userRequests))
	for _, info := range userRequests {
		createInfos = append(createInfos, info)
	}

	createdCharacters := make(map[string]*entity.Character, len(createInfos))
	removeCharNameModels := make([]mongo.WriteModel, 0, len(createInfos))
	newCharModels := make([]mongo.WriteModel, len(createInfos))

	for i, info := range createInfos {
		newCharacter := &entity.Character{
			Slot: info.Slot,
			Name: info.Name,
		}

		update := bson.D{
			{Key: "$push", Value: bson.D{
				{Key: "characters", Value: newCharacter},
			}},
		}
		filter := bson.D{
			{Key: "_id", Value: info.Uid},
		}

		createdCharacters[info.Uid] = newCharacter
		newCharModels[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(false)
	}

	if len(newCharModels) > 0 {
		_, err = r.user.BulkWrite(ctx, newCharModels, options.BulkWrite().SetOrdered(false))
		if err != nil {
			if bulkErr, ok := err.(mongo.BulkWriteException); ok {
				for _, writeErr := range bulkErr.WriteErrors {
					info := createInfos[writeErr.Index]
					delete(createdCharacters, info.Uid)

					opts := mongo.NewDeleteOneModel().SetFilter(bson.D{
						{Key: "_id", Value: info.Name},
					})
					removeCharNameModels = append(removeCharNameModels, opts)
				}
			} else {
				// 캐릭터 생성 중 오류가 발생함... CS처리가 필요 할 수 있음
				// TODO:특수 상황 로그 남기기
				log.Err(err)
				return nil, nil, err
			}
		}
	}

	if len(removeCharNameModels) > 0 {
		_, err = r.characterName.BulkWrite(ctx, removeCharNameModels, options.BulkWrite().SetOrdered(false))
		if err != nil {
			if bulkErr, ok := err.(mongo.BulkWriteException); ok {
				log.Warn().Msg("일부 캐릭터 이름 삭제에 실패했습니다")
				for _, writeErr := range bulkErr.WriteErrors {
					log.Warn().Msgf("캐릭터 이름 삭제 실패: %v - %v", createInfos[writeErr.Index].Name, writeErr.Message)
				}
			}
			log.Err(err).Msg("캐릭터 이름 삭제 중 오류 발생")
		}
	}

	return createdCharacters, failureUids, nil
}

func (r *UserMongoRepository) DeleteCharacters(ctx context.Context, infos []*UserDeleteCharacter) ([]string, error) {
	if len(infos) == 0 {
		return []string{}, nil
	}

	successUids := make(map[string]string, len(infos))
	writeModels := make([]mongo.WriteModel, len(infos))
	for i, info := range infos {
		filter := bson.D{
			{Key: "_id", Value: info.Uid},
		}
		update := bson.D{
			{Key: "$pull", Value: bson.D{
				{Key: "characters", Value: bson.D{
					{Key: "slot", Value: info.Slot},
				}},
			}},
		}

		successUids[info.Uid] = info.Name
		writeModels[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(false)
	}

	_, err := r.user.BulkWrite(ctx, writeModels, options.BulkWrite().SetOrdered(false))
	if err != nil {
		if bulkErr, ok := err.(mongo.BulkWriteException); ok {
			log.Warn().Msg("일부 캐릭터 제거에 실패했습니다")
			for _, writeErr := range bulkErr.WriteErrors {
				log.Warn().Msgf("캐릭터 제거 실패: %v", writeErr.Message)
				delete(successUids, infos[writeErr.Index].Uid)
			}
		} else {
			log.Error().Msg("캐릭터 제거 중 오류 발생")
			return nil, err
		}
	}

	// 캐릭터 삭제에 성공하였다면 캐릭터 이름도 삭제
	if len(successUids) == 0 {
		return []string{}, nil
	}

	deleteNames := make([]string, 0, len(successUids))
	for _, name := range successUids {
		deleteNames = append(deleteNames, name)
	}
	deleteCharNameModels := make([]mongo.WriteModel, 0, len(successUids))

	for _, name := range deleteNames {
		filter := bson.D{
			{Key: "_id", Value: name},
		}

		deleteCharNameModels = append(deleteCharNameModels, mongo.NewDeleteOneModel().SetFilter(filter))
	}

	_, err = r.characterName.BulkWrite(ctx, deleteCharNameModels, options.BulkWrite().SetOrdered(false))
	if err != nil {
		if bulkErr, ok := err.(mongo.BulkWriteException); ok {
			log.Warn().Msg("일부 캐릭터 이름 삭제에 실패했습니다")
			for _, writeErr := range bulkErr.WriteErrors {
				// 삭제에 실패한 캐릭터 이름은 별도로 삭제를 해줘야 함...
				log.Warn().Msgf("캐릭터 이름 삭제 실패: %v - %v", deleteNames[writeErr.Index], writeErr.Message)
			}
		}
		log.Err(err).Msg("캐릭터 이름 삭제 중 오류 발생")
	}
	return func() []string {
		successUidsList := make([]string, 0, len(successUids))
		for uid := range successUids {
			successUidsList = append(successUidsList, uid)
		}
		return successUidsList
	}(), nil
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
