package entity

// 게임 내 랜덤 시드 정보를 담는 엔티티입니다
type GameRandom struct {
	// 캐릭터 생성 시드
	CharacterCreate int64 `bson:"character_create"`
}
