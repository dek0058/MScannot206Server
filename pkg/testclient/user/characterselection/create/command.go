package create

import (
	"MScannot206/pkg/testclient/framework"
	"MScannot206/pkg/testclient/user"
	"MScannot206/pkg/testclient/user/handler"
	"MScannot206/shared/def"
	"errors"
	"strconv"

	"github.com/rs/zerolog/log"
)

func NewCharacterCreateCommand(client framework.Client, userHandler handler.UserHandler) (*CharacterCreateCommand, error) {
	if client == nil {
		return nil, framework.ErrClientIsNil
	}

	if userHandler == nil {
		return nil, handler.ErrUserHandlerIsNil
	}

	userLogic, err := framework.GetLogic[*user.UserLogic](client)
	if err != nil {
		return nil, err
	}

	return &CharacterCreateCommand{
		client:      client,
		userHandler: userHandler,

		userLogic: userLogic,
	}, nil
}

type CharacterCreateCommand struct {
	client      framework.Client
	userHandler handler.UserHandler

	userLogic *user.UserLogic
}

func (c *CharacterCreateCommand) Commands() []string {
	return []string{"character_create"}
}

func (c *CharacterCreateCommand) Execute(args []string) error {
	if len(args) < 2 {
		return framework.ErrInvalidCommandArgument
	}

	slot, err := strconv.Atoi(args[0])
	if err != nil {
		return framework.ErrInvalidCommandArgument
	} else if slot < 0 {
		return errors.New("slot은 0 이상의 값이어야 합니다")
	} else if slot > def.MaxCharacterSlot {
		return errors.New("slot이 최대 캐릭터 슬롯 수를 초과하였습니다")
	}
	name := args[1]

	if err := c.userLogic.RequestCheckCharacterName(c.userHandler.GetUid(), name); err != nil {
		return err
	}

	if err := c.userLogic.RequestCreateCharacter(c.userHandler.GetUid(), slot, name); err != nil {
		return err
	}

	log.Info().Str("Name", name).Msg("캐릭터 생성 요청 완료")

	return nil
}

func (c *CharacterCreateCommand) Description() string {
	return framework.MakeCommandDescription(c.Commands(), "<slot:number> <name>", "캐릭터 생성을 요청 합니다.")
}
