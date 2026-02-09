package entity

import (
	"errors"
	"time"
)

var ErrCharacterSlotIsFull = errors.New("캐릭터 슬롯이 가득 찼습니다")

// 캐릭터 엔티티를 생성 합니다
func NewCharacter(slot int, name string) *Character {
	return &Character{
		Slot: slot,
		Name: name,
	}
}

// 캐릭터 엔티티 구조체
type Character struct {
	// 캐릭터 슬롯 번호
	Slot int `json:"slot" bson:"slot"`

	// 캐릭터 이름
	Name string `json:"name" bson:"name"`
}

var ErrCharacterNameIsNil = errors.New("character name entity is null")

// 캐릭터 이름 엔티티 구조체
type CharacterName struct {
	// 캐릭터 이름
	Name string `bson:"_id"`

	// 생성 일시
	CreatedAt time.Time `bson:"create_at"`
}
