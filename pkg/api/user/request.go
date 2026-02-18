package user

// 캐릭터 이름 검사 요청 정보
type UserNameCheckInfo struct {
	// 유저 고유 ID
	Uid string `json:"uid"`

	// 인증 토큰
	Token string `json:"token"`

	// 검사할 이름
	Name string `json:"name"`
}

// 캐릭터 이름 검사 요청
type CheckCharacterNameRequest struct {
	// 검사 요청 목록
	Requests []*UserNameCheckInfo `json:"requests"`
}

// 캐릭터 생성 요청 정보
type UserCreateCharacterInfo struct {
	// 유저 고유 ID
	Uid string `json:"uid"`

	// 인증 토큰
	Token string `json:"token"`

	// 생성할 캐릭터 슬롯 번호
	Slot int `json:"slot"`

	// 생성할 캐릭터 이름
	Name string `json:"name"`

	// 생성할 캐릭터 성별 (1: 남성, 2: 여성)
	Gender int `json:"gender"`
}

// 캐릭터 생성 요청
type CreateCharacterRequest struct {
	// 생성 요청 목록
	Requests []*UserCreateCharacterInfo `json:"requests"`
}

// 캐릭터 삭제 요청 정보
type UserDeleteCharacterInfo struct {
	// 유저 고유 ID
	Uid string `json:"uid"`

	// 인증 토큰
	Token string `json:"token"`

	// 삭제할 캐릭터 슬롯 번호
	Slot int `json:"slot"`
}

// 캐릭터 삭제 요청
type DeleteCharacterRequest struct {
	// 삭제 요청 목록
	Requests []*UserDeleteCharacterInfo `json:"requests"`
}
