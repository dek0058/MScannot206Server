package login

import (
	"MScannot206/pkg/testclient/framework"
	"MScannot206/shared/service"
	"fmt"
)

func NewLoginCommand(host service.ServiceHost) *LoginCommand {
	return &LoginCommand{
		host: host,
	}
}

type LoginCommand struct {
	host service.ServiceHost
}

func (c *LoginCommand) Commands() []string {
	return []string{"login"}
}

func (c *LoginCommand) Execute(args []string) error {
	loginService, err := service.GetService[*LoginService](c.host)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("파라미터가 부족합니다")
	}

	uid := args[0]

	if err := loginService.RequestLogin(uid); err != nil {
		return err
	}

	return nil
}

func (c *LoginCommand) Description() string {
	return framework.MakeCommandDescription(c.Commands(), "<uid>", "로그인을 요청 합니다.")
}
