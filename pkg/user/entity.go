package user

import (
	"MScannot206/shared/entity"
	"MScannot206/shared/types"
)

// 유저 엔티티 정의
type UserEntity struct {
	// 유저 고유 ID
	Uid string

	// 인증 토큰
	Token string
}

// 캐릭터 생성 정보
type UserCreateCharacter struct {
	// 유저 고유 ID
	Uid string

	// 생성할 캐릭터 슬롯 번호
	Slot int

	// 생성할 캐릭터 이름
	Name string

	// 생성할 캐릭터 성별 (1: 남성, 2: 여성)
	Gender int
}

// 캐릭터 생성 결과
type UserCreateCharacterResult struct {
	// 생성된 캐릭터 정보
	Character *entity.Character

	// 생성된 캐릭터의 장비 정보
	Equips map[types.CharacterEquipType]string

	// 에러 코드
	ErrorCode string
}

// 캐릭터 삭제 정보
type UserDeleteCharacter struct {
	// 유저 고유 ID
	Uid string

	// 삭제할 캐릭터 슬롯 번호
	Slot int

	// 삭제할 캐릭터 이름
	Name string
}
