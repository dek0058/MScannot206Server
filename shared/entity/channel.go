package entity

import "time"

type Channel struct {
	ID        string    `json:"id" bson:"_id"`
	Index     int       `json:"index" bson:"index"`
	LeaseID   string    `json:"lease_id" bson:"lease_id"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
}
