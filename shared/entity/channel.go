package entity

type Channel struct {
	Id    string `json:"id" bson:"_id"`
	Index int    `json:"index" bson:"index"`
}
