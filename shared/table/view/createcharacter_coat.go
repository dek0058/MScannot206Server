package view

import (
	"MScannot206/shared/table"
	"MScannot206/shared/types"
	"MScannot206/shared/util"
	"math/rand/v2"
)

// 캐릭터 생성시 상의/하의/한벌옷을 획득하기 위한 뷰입니다
type CreateCharacterCoatTableView struct {
	createCharacter         *table.CreateCharacterTable
	createCharacterCoat     *table.CreateCharacterCoatTable
	createCharacterPants    *table.CreateCharacterPantsTable
	createCharacterLongCoat *table.CreateCharacterLongCoatTable
}

// 캐릭터 의상을 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterCoatTableView) isHolding(rng *rand.Rand) ([]types.CharacterEquipType, bool) {
	if v.createCharacter == nil {
		return nil, false
	}

	ret := make([]types.CharacterEquipType, 0, 2)
	singleProb := util.NewSingleProbabilisticData[any](1.0, rng)
	var isHoldingLongCoat bool // 2HWeapon 대응: LongCoat

	// 한벌옷(LongCoat)을 먼저 판단합니다. 한벌옷을 획득한 경우 상의(Coat)와 하의(Pants)를 획득할 수 없습니다
	if longCoatRecord, ok := v.createCharacter.Get(string(types.CharacterEquipType_LongCoat)); ok {
		singleProb.Set(types.CharacterEquipType_LongCoat, longCoatRecord.HoldingProb)
		if _, ok := singleProb.Pick(); ok {
			ret = append(ret, types.CharacterEquipType_LongCoat)
			isHoldingLongCoat = true
		}
	}

	// 한벌옷을 획득하지 못한 경우 상의와 하의를 획득할 수 있습니다
	if !isHoldingLongCoat {
		// Coat (1HWeapon 대응)
		if coatRecord, ok := v.createCharacter.Get(string(types.CharacterEquipType_Coat)); ok {
			singleProb.Set(types.CharacterEquipType_Coat, coatRecord.HoldingProb)
			if _, ok := singleProb.Pick(); ok {
				ret = append(ret, types.CharacterEquipType_Coat)
			}
		}

		// Pants (SubWeapon 대응)
		if pantsRecord, ok := v.createCharacter.Get(string(types.CharacterEquipType_Pants)); ok {
			singleProb.Set(types.CharacterEquipType_Pants, pantsRecord.HoldingProb)
			if _, ok := singleProb.Pick(); ok {
				ret = append(ret, types.CharacterEquipType_Pants)
			}
		}
	}

	return ret, len(ret) > 0
}

// 남성 캐릭터 의상을 확률적으로 획득합니다
func (v CreateCharacterCoatTableView) GetMale(rng *rand.Rand) (map[types.CharacterEquipType]string, bool) {
	holdingTypes, ok := v.isHolding(rng)
	if !ok {
		return map[types.CharacterEquipType]string{}, false
	}

	ret := make(map[types.CharacterEquipType]string, len(holdingTypes))

	for _, t := range holdingTypes {
		switch t {
		case types.CharacterEquipType_LongCoat:
			weightedPicker := util.NewWeightedPicker[table.CreateCharacterLongCoatRecord, float64](rng)
			for _, record := range v.GetMaleLongCoatRecords() {
				weightedPicker.Add(record, record.MaleProb)
			}

			pickedRecord, err := weightedPicker.Pick()
			if err == nil {
				ret[types.CharacterEquipType_LongCoat] = pickedRecord.Index
			}

		case types.CharacterEquipType_Coat:
			weightedPicker := util.NewWeightedPicker[table.CreateCharacterCoatRecord, float64](rng)
			for _, record := range v.GetMaleCoatRecords() {
				weightedPicker.Add(record, record.MaleProb)
			}

			pickedRecord, err := weightedPicker.Pick()
			if err == nil {
				ret[types.CharacterEquipType_Coat] = pickedRecord.Index
			}

		case types.CharacterEquipType_Pants:
			weightedPicker := util.NewWeightedPicker[table.CreateCharacterPantsRecord, float64](rng)
			for _, record := range v.GetMalePantsRecords() {
				weightedPicker.Add(record, record.MaleProb)
			}

			pickedRecord, err := weightedPicker.Pick()
			if err == nil {
				ret[types.CharacterEquipType_Pants] = pickedRecord.Index
			}
		}
	}

	return ret, len(ret) > 0
}

// 여성 캐릭터 의상을 확률적으로 획득합니다
func (v CreateCharacterCoatTableView) GetFemale(rng *rand.Rand) (map[types.CharacterEquipType]string, bool) {
	holdingTypes, ok := v.isHolding(rng)
	if !ok {
		return map[types.CharacterEquipType]string{}, false
	}

	ret := make(map[types.CharacterEquipType]string, len(holdingTypes))

	for _, t := range holdingTypes {
		switch t {
		case types.CharacterEquipType_LongCoat:
			weightedPicker := util.NewWeightedPicker[table.CreateCharacterLongCoatRecord, float64](rng)
			for _, record := range v.GetFemaleLongCoatRecords() {
				weightedPicker.Add(record, record.FemaleProb)
			}

			pickedRecord, err := weightedPicker.Pick()
			if err == nil {
				ret[types.CharacterEquipType_LongCoat] = pickedRecord.Index
			}

		case types.CharacterEquipType_Coat:
			weightedPicker := util.NewWeightedPicker[table.CreateCharacterCoatRecord, float64](rng)
			for _, record := range v.GetFemaleCoatRecords() {
				weightedPicker.Add(record, record.FemaleProb)
			}

			pickedRecord, err := weightedPicker.Pick()
			if err == nil {
				ret[types.CharacterEquipType_Coat] = pickedRecord.Index
			}

		case types.CharacterEquipType_Pants:
			weightedPicker := util.NewWeightedPicker[table.CreateCharacterPantsRecord, float64](rng)
			for _, record := range v.GetFemalePantsRecords() {
				weightedPicker.Add(record, record.FemaleProb)
			}

			pickedRecord, err := weightedPicker.Pick()
			if err == nil {
				ret[types.CharacterEquipType_Pants] = pickedRecord.Index
			}
		}
	}

	return ret, len(ret) > 0
}

// 남성 캐릭터 상의 레코드를 모두 가져옵니다
func (v CreateCharacterCoatTableView) GetMaleCoatRecords() []table.CreateCharacterCoatRecord {
	if v.createCharacterCoat == nil {
		return []table.CreateCharacterCoatRecord{}
	}

	records := v.createCharacterCoat.GetAll()
	ret := make([]table.CreateCharacterCoatRecord, 0, len(records))
	for _, record := range records {
		if record.MaleProb > 0 {
			ret = append(ret, record)
		}
	}
	return ret
}

// 여성 캐릭터 상의 레코드를 모두 가져옵니다
func (v CreateCharacterCoatTableView) GetFemaleCoatRecords() []table.CreateCharacterCoatRecord {
	if v.createCharacterCoat == nil {
		return []table.CreateCharacterCoatRecord{}
	}

	records := v.createCharacterCoat.GetAll()
	ret := make([]table.CreateCharacterCoatRecord, 0, len(records))
	for _, record := range records {
		if record.FemaleProb > 0 {
			ret = append(ret, record)
		}
	}
	return ret
}

// 남성 캐릭터 하의 레코드를 모두 가져옵니다
func (v CreateCharacterCoatTableView) GetMalePantsRecords() []table.CreateCharacterPantsRecord {
	if v.createCharacterPants == nil {
		return []table.CreateCharacterPantsRecord{}
	}

	records := v.createCharacterPants.GetAll()
	ret := make([]table.CreateCharacterPantsRecord, 0, len(records))
	for _, record := range records {
		if record.MaleProb > 0 {
			ret = append(ret, record)
		}
	}
	return ret
}

// 여성 캐릭터 하의 레코드를 모두 가져옵니다
func (v CreateCharacterCoatTableView) GetFemalePantsRecords() []table.CreateCharacterPantsRecord {
	if v.createCharacterPants == nil {
		return []table.CreateCharacterPantsRecord{}
	}

	records := v.createCharacterPants.GetAll()
	ret := make([]table.CreateCharacterPantsRecord, 0, len(records))
	for _, record := range records {
		if record.FemaleProb > 0 {
			ret = append(ret, record)
		}
	}
	return ret
}

// 남성 캐릭터 한벌옷 레코드를 모두 가져옵니다
func (v CreateCharacterCoatTableView) GetMaleLongCoatRecords() []table.CreateCharacterLongCoatRecord {
	if v.createCharacterLongCoat == nil {
		return []table.CreateCharacterLongCoatRecord{}
	}

	records := v.createCharacterLongCoat.GetAll()
	ret := make([]table.CreateCharacterLongCoatRecord, 0, len(records))
	for _, record := range records {
		if record.MaleProb > 0 {
			ret = append(ret, record)
		}
	}
	return ret
}

// 여성 캐릭터 한벌옷 레코드를 모두 가져옵니다
func (v CreateCharacterCoatTableView) GetFemaleLongCoatRecords() []table.CreateCharacterLongCoatRecord {
	if v.createCharacterLongCoat == nil {
		return []table.CreateCharacterLongCoatRecord{}
	}

	records := v.createCharacterLongCoat.GetAll()
	ret := make([]table.CreateCharacterLongCoatRecord, 0, len(records))
	for _, record := range records {
		if record.FemaleProb > 0 {
			ret = append(ret, record)
		}
	}
	return ret
}
