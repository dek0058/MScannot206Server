package user

import (
	user_api "MScannot206/pkg/api/user"
	"MScannot206/pkg/testclient/framework"
	"MScannot206/pkg/testclient/user/character"
	"MScannot206/pkg/user"
	"MScannot206/shared"
	"MScannot206/shared/entity"
	"errors"
)

const userCapacity = 1000

var ErrUserNotFound = errors.New("유저를 찾지 못하였습니다")

func NewUserLogic(client framework.Client) (*UserLogic, error) {
	if client == nil {
		return nil, framework.ErrClientIsNil
	}
	return &UserLogic{
		client: client,

		users: make(map[string]*User, userCapacity),
	}, nil
}

type UserLogic struct {
	client framework.Client

	users map[string]*User
}

func (l *UserLogic) Init() error {
	return nil
}

func (l *UserLogic) Start() error {
	return nil
}

func (l *UserLogic) Stop() error {
	return nil
}

func (l *UserLogic) ConnectUser(userEntity *entity.User, token string) (*User, error) {
	var errs error

	if userEntity == nil {
		return nil, entity.ErrUserIsNil
	}

	u, err := NewUser(userEntity.Uid, token)
	if err != nil {
		return nil, err
	}
	u.Characters = make([]*character.Character, 0, len(userEntity.Characters))
	for _, entry := range userEntity.Characters {
		ch, err := character.NewCharacter(entry.Slot, entry.Name)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		u.Characters = append(u.Characters, ch)
	}

	if errs != nil {
		return nil, errs
	}

	l.users[userEntity.Uid] = u
	return u, nil
}

func (l *UserLogic) DisconnectUser(uid string) error {
	if u, ok := l.users[uid]; ok {
		u.Quit()
		delete(l.users, uid)
	}
	return nil
}

func (l *UserLogic) GetUser(uid string) (*User, bool) {
	user, ok := l.users[uid]
	return user, ok
}

func (l *UserLogic) GetCharacterSlotCount(uid string) (int, error) {
	u, ok := l.users[uid]
	if !ok {
		return 0, ErrUserNotFound
	}
	var slotCount int
	for _, ch := range u.Characters {
		if ch != nil {
			slotCount++
		}
	}
	return slotCount, nil
}

func (l *UserLogic) RequestCheckCharacterName(uid string, name string) error {
	u, ok := l.users[uid]
	if !ok {
		return ErrUserNotFound
	}

	req := &user_api.CheckCharacterNameRequest{
		Requests: []*user_api.UserNameCheckInfo{
			{
				Uid:   u.Uid,
				Token: u.Token,
				Name:  name,
			},
		},
	}

	res, err := framework.WebRequest[user_api.CheckCharacterNameRequest, user_api.CheckCharacterNameResponse](l.client).
		Endpoint("api/v1/user/character/create/check_name").
		Body(req).
		Post()

	if err != nil {
		return err
	}

	if len(res.Responses) == 0 {
		return shared.ToError(user.USER_CHECK_CHARACTER_NAME_UNKNOWN_ERROR)
	}

	var errorCode string = ""
	for _, r := range res.Responses {
		if r.Uid == uid {
			errorCode = r.ErrorCode
			break
		}
	}

	if errorCode != "" {
		return shared.ToError(errorCode)
	}

	return nil
}

func (l *UserLogic) RequestCreateCharacter(uid string, slot int, name string) error {
	u, ok := l.users[uid]
	if !ok {
		return ErrUserNotFound
	}

	req := &user_api.CreateCharacterRequest{
		Requests: []*user_api.UserCreateCharacterInfo{
			{
				Uid:   u.Uid,
				Token: u.Token,
				Slot:  slot,
				Name:  name,
			},
		},
	}

	res, err := framework.WebRequest[user_api.CreateCharacterRequest, user_api.CreateCharacterResponse](l.client).
		Endpoint("api/v1/user/character/create").
		Body(req).
		Post()

	if err != nil {
		return err
	}

	if len(res.Responses) == 0 {
		return shared.ToError(user.USER_CREATE_CHARACTER_UNKNOWN_ERROR)
	}

	var response *user_api.UserCreateCharacterResult
	for _, r := range res.Responses {
		if r.Uid != uid {
			continue
		}
		response = r
	}

	if response.ErrorCode != "" {
		return shared.ToError(response.ErrorCode)
	}

	if response.Character == nil {
		return shared.ToError(user.USER_CREATE_CHARACTER_UNKNOWN_ERROR)
	}

	newCh, err := character.NewCharacter(response.Character.Slot, response.Character.Name)
	if err != nil {
		return err
	}

	u.Characters = append(u.Characters, newCh)

	return nil
}

func (l *UserLogic) RequestDeleteCharacter(uid string, slot int) error {
	u, ok := l.users[uid]
	if !ok {
		return ErrUserNotFound
	}

	req := &user_api.DeleteCharacterRequest{
		Requests: []*user_api.UserDeleteCharacterInfo{
			{
				Uid:   u.Uid,
				Token: u.Token,
				Slot:  slot,
			},
		},
	}

	res, err := framework.WebRequest[user_api.DeleteCharacterRequest, user_api.DeleteCharacterResponse](l.client).
		Endpoint("api/v1/user/character/delete").
		Body(req).
		Post()

	if err != nil {
		return err
	}

	if len(res.Responses) == 0 {
		return shared.ToError(user.USER_DELETE_CHARACTER_UNKNOWN_ERROR)
	}

	var response *user_api.UserDeleteCharacterResult
	for _, r := range res.Responses {
		if r.Uid == uid {
			response = r
			break
		}
	}

	if response.ErrorCode != "" {
		return shared.ToError(response.ErrorCode)
	}

	for i, ch := range u.Characters {
		if ch.Slot == slot {
			u.Characters = append(u.Characters[:i], u.Characters[i+1:]...)
			break
		}
	}

	return nil
}
