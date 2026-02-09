package user

import "MScannot206/shared/entity"

// 캐릭터 이름 검사 결과
type UserNameCheckResult struct {
	// 유저 고유 ID
	Uid string `json:"uid"`

	// 이름 사용 가능 여부
	ErrorCode string `json:"error_code,omitempty"`
}

// 캐릭터 이름 검사 응답
type CheckCharacterNameResponse struct {
	// 검사 결과 목록
	Responses []*UserNameCheckResult `json:"responses"`
}

// 캐릭터 생성 결과
type UserCreateCharacterResult struct {
	// 유저 고유 ID
	Uid string `json:"uid"`

	// 생성된 캐릭터 정보
	Character *entity.Character `json:"character,omitempty"`

	// 생성 오류 코드
	ErrorCode string `json:"error_code,omitempty"`
}

// 캐릭터 생성 응답
type CreateCharacterResponse struct {
	Responses []*UserCreateCharacterResult `json:"responses"`
}

// 캐릭터 삭제 결과
type UserDeleteCharacterResult struct {
	// 유저 고유 ID
	Uid string `json:"uid"`

	// 삭제된 캐릭터 슬롯 번호
	Slot int `json:"slot"`

	// 삭제 오류 코드
	ErrorCode string `json:"error_code,omitempty"`
}

// 캐릭터 삭제 응답
type DeleteCharacterResponse struct {
	// 삭제 결과 목록
	Responses []*UserDeleteCharacterResult `json:"responses"`
}
