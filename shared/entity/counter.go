package entity

type Counter struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}
