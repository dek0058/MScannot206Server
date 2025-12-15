package entity

type UserSession struct {
	Uid       string `bson:"_id"`
	Token     string `bson:"access_token"`
	UpdatedAt int64  `bson:"updated_at"`
}
