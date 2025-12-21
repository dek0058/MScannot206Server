package entity

import "errors"

var ErrUserIsNil = errors.New("user entity is nil")

func NewUser(uid string) *User {
	return &User{
		Uid:        uid,
		Characters: []*Character{},
	}
}

type User struct {
	Uid string `json:"uid" bson:"_id"`

	Characters []*Character `json:"characters,omitempty" bson:"characters,omitempty"`
}
