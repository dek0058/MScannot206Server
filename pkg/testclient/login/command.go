package login

import (
	"MScannot206/pkg/testclient/framework"
	"fmt"
)

func NewLoginCommand(client framework.Client) (*LoginCommand, error) {
	return &LoginCommand{
		client: client,
	}, nil
}

type LoginCommand struct {
	client framework.Client
}

func (c *LoginCommand) Commands() []string {
	return []string{"login"}
}

func (c *LoginCommand) Execute(args []string) error {
	loginLogic, err := framework.GetLogic[*LoginLogic](c.client)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("파라미터가 부족합니다")
	}

	uid := args[0]

	if err := loginLogic.RequestLogin(uid); err != nil {
		return err
	}

	return nil
}

func (c *LoginCommand) Description() string {
	return framework.MakeCommandDescription(c.Commands(), "<uid>", "로그인을 요청 합니다.")
}
