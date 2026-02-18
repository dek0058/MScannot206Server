package entity

import "MScannot206/shared/types"

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

	// 캐릭터 성별
	Gender int `json:"gender" bson:"gender"`
}

// 캐릭터 이름 엔티티 구조체
type CharacterName struct {
	// 캐릭터 이름
	Name string `bson:"_id"`

	// 생성 일시
	CreatedAt int64 `bson:"created_at"`
}

// 캐릭터 장비 엔티티 구조체
type CharacterEquip struct {
	// 장비 종류
	Type types.CharacterEquipType `json:"type" bson:"type"`

	// 장비 인덱스
	Index string `json:"index" bson:"index"`
}

// 캐릭터 장비 슬롯 엔티티 구조체
type CharacterEquipSlot struct {
	// 장비 슬롯 타입
	Type types.CharacterEquipType `json:"type" bson:"type"`

	// 장비 아이템 Id
	ItemId int `json:"item_id" bson:"item_id"`
}
