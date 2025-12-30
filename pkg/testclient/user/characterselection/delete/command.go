package delete

import (
	"MScannot206/pkg/testclient/framework"
	"MScannot206/pkg/testclient/user"
	"MScannot206/pkg/testclient/user/handler"
	"MScannot206/shared/def"
	"errors"
	"strconv"

	"github.com/rs/zerolog/log"
)

func NewCharacterDeleteCommand(client framework.Client, userHandler handler.UserHandler) (*CharacterDeleteCommand, error) {
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

	return &CharacterDeleteCommand{
		client:      client,
		userHandler: userHandler,

		userLogic: userLogic,
	}, nil
}

type CharacterDeleteCommand struct {
	client      framework.Client
	userHandler handler.UserHandler

	userLogic *user.UserLogic
}

func (c *CharacterDeleteCommand) Commands() []string {
	return []string{"character_delete"}
}

func (c *CharacterDeleteCommand) Execute(args []string) error {
	if len(args) < 1 {
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

	if err := c.userLogic.RequestDeleteCharacter(c.userHandler.GetUid(), slot); err != nil {
		return err
	}

	log.Info().Int("Slot", slot).Msg("캐릭터 삭제를 요청 완료")

	return nil
}

func (c *CharacterDeleteCommand) Description() string {
	return framework.MakeCommandDescription(c.Commands(), "<slot:number>", "캐릭터 삭제를 요청 합니다.")
}
