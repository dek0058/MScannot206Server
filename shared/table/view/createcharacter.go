package view

import (
	"MScannot206/shared/table"
	"MScannot206/shared/types"
	"math/rand/v2"
)

func NewCreateCharacterView(
	createCharacter table.CreateCharacterTable,
	createCharacterHair table.CreateCharacterHairTable,
	createCharacterFace table.CreateCharacterFaceTable,
	createCharacterCap table.CreateCharacterCapTable,
	createCharacterCape table.CreateCharacterCapeTable,
	createCharacterCoat table.CreateCharacterCoatTable,
	createCharacterGlove table.CreateCharacterGloveTable,
	createCharacterLongCoat table.CreateCharacterLongCoatTable,
	createCharacterPants table.CreateCharacterPantsTable,
	createCharacterShoes table.CreateCharacterShoesTable,
	createCharacterFaceAccessory table.CreateCharacterFaceAccTable,
	createCharacterEyeAccessory table.CreateCharacterEysAccTable,
	createCharacterEarAccessory table.CreateCharacterEarAccTable,
	createCharacter1HWeapon table.CreateCharacter1HWeaponTable,
	createCharacter2HWeapon table.CreateCharacter2HWeaponTable,
	createCharacterSubWeapon table.CreateCharacterSubWeaponTable,
	createCharacterEar table.CreateCharacterEarTable,
	createCharacterSkin table.CreateCharacterSkinTable,
) CreateCharacterView {
	return CreateCharacterView{
		HairView:          CreateCharacterHairTableView{&createCharacter, &createCharacterHair},
		FaceView:          CreateCharacterFaceTableView{&createCharacter, &createCharacterFace},
		CapView:           CreateCharacterCapTableView{&createCharacter, &createCharacterCap},
		CapeView:          CreateCharacterCapeTableView{&createCharacter, &createCharacterCape},
		CoatView:          CreateCharacterCoatTableView{&createCharacter, &createCharacterCoat, &createCharacterPants, &createCharacterLongCoat},
		GloveView:         CreateCharacterGloveTableView{&createCharacter, &createCharacterGlove},
		ShoesView:         CreateCharacterShoesTableView{&createCharacter, &createCharacterShoes},
		FaceAccessoryView: CreateCharacterFaceAccTableView{&createCharacter, &createCharacterFaceAccessory},
		EyeAccessoryView:  CreateCharacterEysAccTableView{&createCharacter, &createCharacterEyeAccessory},
		EarAccessoryView:  CreateCharacterEarAccTableView{&createCharacter, &createCharacterEarAccessory},
		WeaponView:        CreateCharacterWeaponTableView{&createCharacter, &createCharacter1HWeapon, &createCharacter2HWeapon, &createCharacterSubWeapon},
		EarView:           CreateCharacterEarTableView{&createCharacter, &createCharacterEar},
		SkinView:          CreateCharacterSkinTableView{&createCharacter, &createCharacterSkin},
	}
}

type CreateCharacterView struct {
	HairView          CreateCharacterHairTableView
	FaceView          CreateCharacterFaceTableView
	CapView           CreateCharacterCapTableView
	CapeView          CreateCharacterCapeTableView
	CoatView          CreateCharacterCoatTableView
	GloveView         CreateCharacterGloveTableView
	ShoesView         CreateCharacterShoesTableView
	FaceAccessoryView CreateCharacterFaceAccTableView
	EyeAccessoryView  CreateCharacterEysAccTableView
	EarAccessoryView  CreateCharacterEarAccTableView
	WeaponView        CreateCharacterWeaponTableView
	EarView           CreateCharacterEarTableView
	SkinView          CreateCharacterSkinTableView
}

func (v CreateCharacterView) GetMale(rng *rand.Rand) map[types.CharacterEquipType]string {
	var ret = make(map[types.CharacterEquipType]string, types.CharacterEquipCount)

	if hair, ok := v.HairView.GetMale(rng); ok {
		ret[types.CharacterEquipType_Hair] = hair
	}

	if face, ok := v.FaceView.GetMale(rng); ok {
		ret[types.CharacterEquipType_Face] = face
	}

	if cap, ok := v.CapView.GetMale(rng); ok {
		ret[types.CharacterEquipType_Cap] = cap
	}

	if cape, ok := v.CapeView.GetMale(rng); ok {
		ret[types.CharacterEquipType_Cape] = cape
	}

	if clots, ok := v.CoatView.GetMale(rng); ok {
		for equipType, index := range clots {
			ret[equipType] = index
		}
	}

	if glove, ok := v.GloveView.GetMale(rng); ok {
		ret[types.CharacterEquipType_Glove] = glove
	}

	if shoes, ok := v.ShoesView.GetMale(rng); ok {
		ret[types.CharacterEquipType_Shoes] = shoes
	}

	if faceAcc, ok := v.FaceAccessoryView.GetMale(rng); ok {
		ret[types.CharacterEquipType_FaceAccessory] = faceAcc
	}

	if eyeAcc, ok := v.EyeAccessoryView.GetMale(rng); ok {
		ret[types.CharacterEquipType_EyeAccessory] = eyeAcc
	}

	if earAcc, ok := v.EarAccessoryView.GetMale(rng); ok {
		ret[types.CharacterEquipType_EarAccessory] = earAcc
	}

	if weapons, ok := v.WeaponView.GetMale(rng); ok {
		for equipType, index := range weapons {
			ret[equipType] = index
		}
	}

	if ear, ok := v.EarView.GetMale(rng); ok {
		ret[types.CharacterEquipType_Ear] = ear
	}

	if skin, ok := v.SkinView.GetMale(rng); ok {
		ret[types.CharacterEquipType_Skin] = skin
	}

	return ret
}

func (v CreateCharacterView) GetFemale(rng *rand.Rand) map[types.CharacterEquipType]string {
	var ret = make(map[types.CharacterEquipType]string, types.CharacterEquipCount)

	if hair, ok := v.HairView.GetFemale(rng); ok {
		ret[types.CharacterEquipType_Hair] = hair
	}

	if face, ok := v.FaceView.GetFemale(rng); ok {
		ret[types.CharacterEquipType_Face] = face
	}

	if cap, ok := v.CapView.GetFemale(rng); ok {
		ret[types.CharacterEquipType_Cap] = cap
	}

	if cape, ok := v.CapeView.GetFemale(rng); ok {
		ret[types.CharacterEquipType_Cape] = cape
	}

	if clots, ok := v.CoatView.GetFemale(rng); ok {
		for equipType, index := range clots {
			ret[equipType] = index
		}
	}

	if glove, ok := v.GloveView.GetFemale(rng); ok {
		ret[types.CharacterEquipType_Glove] = glove
	}

	if shoes, ok := v.ShoesView.GetFemale(rng); ok {
		ret[types.CharacterEquipType_Shoes] = shoes
	}

	if faceAcc, ok := v.FaceAccessoryView.GetFemale(rng); ok {
		ret[types.CharacterEquipType_FaceAccessory] = faceAcc
	}

	if eyeAcc, ok := v.EyeAccessoryView.GetFemale(rng); ok {
		ret[types.CharacterEquipType_EyeAccessory] = eyeAcc
	}

	if earAcc, ok := v.EarAccessoryView.GetFemale(rng); ok {
		ret[types.CharacterEquipType_EarAccessory] = earAcc
	}

	if weapons, ok := v.WeaponView.GetFemale(rng); ok {
		for equipType, index := range weapons {
			ret[equipType] = index
		}
	}

	if ear, ok := v.EarView.GetFemale(rng); ok {
		ret[types.CharacterEquipType_Ear] = ear
	}

	if skin, ok := v.SkinView.GetFemale(rng); ok {
		ret[types.CharacterEquipType_Skin] = skin
	}

	return ret
}
