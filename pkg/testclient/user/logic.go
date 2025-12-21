package user

import (
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

	req := &user.CheckCharacterNameRequest{
		Requests: []*user.UserNameCheckInfo{
			{
				Uid:   u.Uid,
				Token: u.Token,
				Name:  name,
			},
		},
	}

	res, err := framework.WebRequest[user.CheckCharacterNameRequest, user.CheckCharacterNameResponse](l.client).
		Endpoint("user/character/create/check_name").
		Body(req).
		Post()

	if err != nil {
		return err
	}

	if len(res.Responses) == 0 {
		return shared.ToError(user.USER_CHECK_CHARACTER_NAME_UNKNOWN_ERROR)
	}

	var available bool = false
	var errorCode string = ""
	for _, resp := range res.Responses {
		if resp.Uid == uid {
			available = resp.Available
			errorCode = resp.ErrorCode
			break
		}
	}

	if !available {
		return shared.ToError(errorCode)
	}

	return nil
}

func (l *UserLogic) RequestCreateCharacter(uid string, slot int, name string) error {
	u, ok := l.users[uid]
	if !ok {
		return ErrUserNotFound
	}

	req := &user.CreateCharacterRequest{
		Requests: []*user.UserCreateCharacterInfo{
			{
				Uid:   u.Uid,
				Token: u.Token,
				Slot:  slot,
				Name:  name,
			},
		},
	}

	res, err := framework.WebRequest[user.CreateCharacterRequest, user.CreateCharacterResponse](l.client).
		Endpoint("user/character/create").
		Body(req).
		Post()

	if err != nil {
		return err
	}

	if len(res.Responses) == 0 {
		return shared.ToError(user.USER_CREATE_CHARACTER_UNKNOWN_ERROR)
	}

	var errorCode string = user.USER_CREATE_CHARACTER_UNKNOWN_ERROR
	for _, resp := range res.Responses {
		if resp.Uid == uid && resp.Slot == slot {
			errorCode = resp.ErrorCode
			break
		}
	}

	if errorCode != "" {
		return shared.ToError(errorCode)
	}

	newCh, err := character.NewCharacter(slot, name)
	if err != nil {
		return err
	}

	u.Characters = append(u.Characters, newCh)

	return nil
}
