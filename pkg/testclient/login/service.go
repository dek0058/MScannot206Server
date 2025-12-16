package login

import (
	"MScannot206/pkg/login"
	"MScannot206/shared/client"
	"MScannot206/shared/entity"
	"errors"
	"fmt"
)

func NewLoginService(client *client.WebClient) (*LoginService, error) {
	if client == nil {
		return nil, errors.New("http client is nil")
	}

	s := &LoginService{
		client: client,

		users: make(map[string]*entity.User),
	}

	return s, nil
}

type LoginService struct {
	client *client.WebClient

	users map[string]*entity.User
}

func (s *LoginService) Init() error {
	return nil
}

func (s *LoginService) Start() error {
	return nil
}

func (s *LoginService) Stop() error {
	return nil
}

func (s *LoginService) GetUser(uid string) (*entity.User, error) {
	if user, ok := s.users[uid]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("%v is not found", uid)
}

func (s *LoginService) LoginUser(uid string) {
	user := entity.NewUser(uid)
	s.users[uid] = user
}

func (s *LoginService) RequestLogin(uid string) error {
	if uid == "" {
		return fmt.Errorf("uid is empty")
	}

	if GetUser, ok := s.users[uid]; ok {
		return fmt.Errorf("이미 로그인된 사용자입니다: %s", GetUser.Uid)
	}

	req := &login.LoginRequest{
		Uids: []string{uid},
	}

	fmt.Printf("로그인 요청: %s\n", uid)

	res, err := client.WebRequest[login.LoginRequest, login.LoginResponse](s.client).Endpoint("login").Body(req).Post()
	if err != nil {
		return err
	}

	_ = res

	// test

	// var success bool = false
	// for _, successUid := range res.SuccessUids {
	// 	if successUid == uid {
	// 		success = true
	// 		break
	// 	}
	// }

	// if !success {
	// 	for _, failUid := range res.FailUids {
	// 		if failUid.Uid == uid {
	// 			return manager.ToError(failUid.ErrorCode)
	// 		}
	// 	}

	// 	return manager.ToError(login.LOGIN_LOGIN_UNABLE)
	// }

	// fmt.Printf("로그인 성공: %s\n", uid)

	return nil
}
