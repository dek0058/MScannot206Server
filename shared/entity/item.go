package entity

// 아이템 엔티티를 생성 합니다
func NewItem(id, index string, count int64, bound bool) *Item {
	return &Item{
		Id:    id,
		Index: index,
		Count: count,
		Bound: bound,
	}
}

// 아이템 엔티티 구조체
type Item struct {
	// 아이템 ID
	Id string `json:"id" bson:"_id"`

	// 아이템 테이블 인덱스
	Index string `json:"index" bson:"index"`

	// 아이템 개수
	Count int64 `json:"count" bson:"count"`

	// 아이템이 귀속되어 있는지 여부
	Bound bool `json:"bound" bson:"bound"`
}
