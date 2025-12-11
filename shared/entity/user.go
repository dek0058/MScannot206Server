package entity

func NewUser(uid string) *User {
	return &User{
		Uid: uid,
	}
}

type User struct {
	Uid string `json:"uid" bson:"uid"`
}
