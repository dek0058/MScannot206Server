package entity

import "errors"

var ErrUserIsNil = errors.New("user entity is nil")

// User 엔티티를 생성 합니다
func NewUser(uid string) *User {
	return &User{
		Uid:        uid,
		Characters: []*Character{},
	}
}

// 유저 엔티티 구조체
type User struct {
	// 유저 고유 아이디
	Uid string `json:"uid" bson:"_id"`

	// 유저가 보유한 캐릭터 목록
	Characters []*Character `json:"characters,omitempty" bson:"characters,omitempty"`
}
