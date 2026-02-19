package view_test

import (
	"MScannot206/shared/table"
	"MScannot206/shared/table/view"
	"MScannot206/shared/types"
	"maps"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateCharacterView(t *testing.T) {

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	dataPath := filepath.Join(wd, "../../../data")

	tableRepo := &table.Repository{}
	if err := tableRepo.Load(dataPath); err != nil {
		t.Fatalf("failed to load table repository: %v", err)
	}

	view := view.NewCreateCharacterView(
		tableRepo.CreateCharacter,
		tableRepo.CreateCharacterHair,
		tableRepo.CreateCharacterFace,
		tableRepo.CreateCharacterCap,
		tableRepo.CreateCharacterCape,
		tableRepo.CreateCharacterCoat,
		tableRepo.CreateCharacterGlove,
		tableRepo.CreateCharacterLongCoat,
		tableRepo.CreateCharacterPants,
		tableRepo.CreateCharacterShoes,
		tableRepo.CreateCharacterFaceAcc,
		tableRepo.CreateCharacterEysAcc,
		tableRepo.CreateCharacterEarAcc,
		tableRepo.CreateCharacter1HWeapon,
		tableRepo.CreateCharacter2HWeapon,
		tableRepo.CreateCharacterSubWeapon,
		tableRepo.CreateCharacterEar,
		tableRepo.CreateCharacterSkin,
	)

	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewPCG(uint64(seed), uint64(seed)))

	testCases := []struct {
		gender int
	}{
		{gender: types.GenderType_Male},
		{gender: types.GenderType_Female},
	}

	for _, tc := range testCases {
		t.Run("CreateCharacterView", func(t *testing.T) {
			ret := map[types.CharacterEquipType]string{}
			switch tc.gender {
			case types.GenderType_Male:
				t.Logf("Gender: Male")
				maps.Copy(ret, view.GetMale(rng))
			case types.GenderType_Female:
				t.Logf("Gender: Female")
				maps.Copy(ret, view.GetFemale(rng))
			}

			for equipType, index := range ret {
				t.Logf("equipType: %v, index: %v", equipType, index)
			}
		})
	}
}
