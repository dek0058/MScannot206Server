package entity

import "time"

type Channel struct {
	Id        string    `json:"id" bson:"_id"`
	Index     int       `json:"index" bson:"index"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
}

type ChannelRecycler struct {
	Id    string `json:"id" bson:"_id"`
	Index int    `json:"index" bson:"index"`
}
