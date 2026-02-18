package user

import (
	"errors"
	"math/rand/v2"
)

var ErrRandomServiceHandlerIsNil = errors.New("random service handler is null")

// 랜덤 서비스 핸들러는 유저 서비스에서 랜덤 시드 정보를 얻기 위해 사용하는 핸들러입니다
type RandomServiceHandler interface {
	GetCharacterCreateSeed() *rand.Rand
}
