package login

import (
	"MScannot206/pkg/login"
	"MScannot206/pkg/testclient/framework"
	"MScannot206/shared"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

func NewLoginLogic(client framework.Client) (*LoginLogic, error) {
	if client == nil {
		return nil, errors.New("http client is nil")
	}

	s := &LoginLogic{
		client: client,
	}

	return s, nil
}

type LoginLogic struct {
	client framework.Client
}

func (l *LoginLogic) Init() error {
	return nil
}

func (l *LoginLogic) Start() error {
	return nil
}

func (l *LoginLogic) Stop() error {
	return nil
}
func (l *LoginLogic) RequestLogin(uid string) error {
	if uid == "" {
		return fmt.Errorf("uid is empty")
	}

	req := &login.LoginRequest{
		Uids: []string{uid},
	}

	log.Info().Msgf("로그인 요청: %s", uid)

	res, err := framework.WebRequest[login.LoginRequest, login.LoginResponse](l.client).Endpoint("login").Body(req).Post()
	if err != nil {
		return err
	}

	successCount := len(res.SuccessUids)
	failCount := len(res.FailUids)

	token := ""

	if successCount == 0 && failCount == 0 {
		return shared.ToError(login.LOGIN_LOGIN_UNABLE)
	} else if failCount > 0 {
		for _, failUid := range res.FailUids {
			if failUid.Uid == uid {
				return shared.ToError(failUid.ErrorCode)
			}
		}
	} else if successCount > 0 {
		for _, successUid := range res.SuccessUids {
			if successUid.Uid == uid {
				token = successUid.Token
				break
			}
		}
	}

	if token == "" {
		return shared.ToError(login.LOGIN_UNKOWN_ERROR)
	}

	log.Info().Msgf("로그인 성공: %s, 토큰: %s", uid, token)

	return nil
}
