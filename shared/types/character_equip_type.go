package types

// CharacterEquipType은 캐릭터의 장비 종류를 나타내는 타입입니다
type CharacterEquipType string

// 캐릭터 장비 타입 개수
const CharacterEquipCount = 17

// 캐릭터 장비 슬롯 개수
const CharacterEquipSlotCount = 15

const (
	CharacterEquipType_None = CharacterEquipType("none")

	// 헤어
	CharacterEquipType_Hair = CharacterEquipType("hair")

	// 얼굴
	CharacterEquipType_Face = CharacterEquipType("face")

	// 모자
	CharacterEquipType_Cap = CharacterEquipType("cap")

	// 망토
	CharacterEquipType_Cape = CharacterEquipType("cape")

	// 상의
	CharacterEquipType_Coat = CharacterEquipType("coat")

	// 장갑
	CharacterEquipType_Glove = CharacterEquipType("glove")

	// 한벌 옷
	CharacterEquipType_LongCoat = CharacterEquipType("longcoat")

	// 하의
	CharacterEquipType_Pants = CharacterEquipType("pants")

	// 신발
	CharacterEquipType_Shoes = CharacterEquipType("shoes")

	// 얼굴 장식
	CharacterEquipType_FaceAccessory = CharacterEquipType("faceaccessory")

	// 눈 장식
	CharacterEquipType_EyeAccessory = CharacterEquipType("eyeaccessory")

	// 귀 장식
	CharacterEquipType_EarAccessory = CharacterEquipType("earaccessory")

	// 한손 무기
	CharacterEquipType_1HWeapon = CharacterEquipType("onehandedweapon")

	// 두손 무기
	CharacterEquipType_2HWeapon = CharacterEquipType("twohandedweapon")

	// 보조 무기
	CharacterEquipType_SubWeapon = CharacterEquipType("subweapon")

	// 귀
	CharacterEquipType_Ear = CharacterEquipType("ear")

	// 피부
	CharacterEquipType_Skin = CharacterEquipType("skin")
)
