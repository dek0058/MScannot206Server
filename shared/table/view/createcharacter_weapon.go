package view

import (
	"MScannot206/shared/table"
	"MScannot206/shared/types"
	"MScannot206/shared/util"
	"math/rand/v2"
)

// 캐릭터 생성시 무기를 획득하기 위한 뷰입니다
type CreateCharacterWeaponTableView struct {
	createCharacter          *table.CreateCharacterTable
	createCharacter1HWeapon  *table.CreateCharacter1HWeaponTable
	createCharacter2HWeapon  *table.CreateCharacter2HWeaponTable
	createCharacterSubWeapon *table.CreateCharacterSubWeaponTable
}

// 캐릭터 무기를 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterWeaponTableView) isHolding(rng *rand.Rand) ([]types.CharacterEquipType, bool) {
	if v.createCharacter == nil {
		return nil, false
	}

	ret := make([]types.CharacterEquipType, 0, 2)
	singleProb := util.NewSingleProbabilisticData[any](1.0, rng)
	var isHolding2HWeapon bool

	// 2H 무기를 먼저 판단합니다. 2H 무기를 획득한 경우 1H 무기와 보조 무기를 획득할 수 없습니다
	if twoHandRecord, ok := v.createCharacter.Get(string(types.CharacterEquipType_2HWeapon)); ok {
		singleProb.Set(types.CharacterEquipType_2HWeapon, twoHandRecord.HoldingProb)
		if _, ok := singleProb.Pick(); ok {
			ret = append(ret, types.CharacterEquipType_2HWeapon)
			isHolding2HWeapon = true
		}
	}

	// 2H 무기를 획득한 경우 1H 무기와 보조 무기를 획득할 수 없습니다
	if !isHolding2HWeapon {
		if oneHandRecord, ok := v.createCharacter.Get(string(types.CharacterEquipType_1HWeapon)); ok {
			singleProb.Set(types.CharacterEquipType_1HWeapon, oneHandRecord.HoldingProb)
			if _, ok := singleProb.Pick(); ok {
				ret = append(ret, types.CharacterEquipType_1HWeapon)
			}
		}

		if subWeaponRecord, ok := v.createCharacter.Get(string(types.CharacterEquipType_SubWeapon)); ok {
			singleProb.Set(types.CharacterEquipType_SubWeapon, subWeaponRecord.HoldingProb)
			if _, ok := singleProb.Pick(); ok {
				ret = append(ret, types.CharacterEquipType_SubWeapon)
			}
		}
	}

	return ret, len(ret) > 0
}

// 남성 캐릭터 무기를 확률적으로 획득합니다
func (v CreateCharacterWeaponTableView) GetMale(rng *rand.Rand) (map[types.CharacterEquipType]string, bool) {
	holdingTypes, ok := v.isHolding(rng)
	if !ok {
		return map[types.CharacterEquipType]string{}, false
	}

	ret := make(map[types.CharacterEquipType]string, len(holdingTypes))

	for _, t := range holdingTypes {
		switch t {
		case types.CharacterEquipType_1HWeapon:
			weightedPicker := util.NewWeightedPicker[table.CreateCharacter1HWeaponRecord, float64](rng)
			for _, record := range v.GetMale1HWeapon(rng) {
				weightedPicker.Add(record, record.MaleProb)
			}

			pickedRecord, err := weightedPicker.Pick()
			if err == nil {
				ret[types.CharacterEquipType_1HWeapon] = pickedRecord.Index
			}

		case types.CharacterEquipType_2HWeapon:
			weightedPicker := util.NewWeightedPicker[table.CreateCharacter2HWeaponRecord, float64](rng)
			for _, record := range v.GetMale2HWeapon(rng) {
				weightedPicker.Add(record, record.MaleProb)
			}

			pickedRecord, err := weightedPicker.Pick()
			if err == nil {
				ret[types.CharacterEquipType_2HWeapon] = pickedRecord.Index
			}

		case types.CharacterEquipType_SubWeapon:
			weightedPicker := util.NewWeightedPicker[table.CreateCharacterSubWeaponRecord, float64](rng)
			for _, record := range v.GetMaleSubWeapon(rng) {
				weightedPicker.Add(record, record.MaleProb)
			}

			pickedRecord, err := weightedPicker.Pick()
			if err == nil {
				ret[types.CharacterEquipType(pickedRecord.SubType)] = pickedRecord.Index
			}
		}
	}

	return ret, len(ret) > 0
}

// 여성 캐릭터 무기를 확률적으로 획득합니다
func (v CreateCharacterWeaponTableView) GetFemale(rng *rand.Rand) (map[types.CharacterEquipType]string, bool) {
	holdingTypes, ok := v.isHolding(rng)
	if !ok {
		return map[types.CharacterEquipType]string{}, false
	}

	ret := make(map[types.CharacterEquipType]string, len(holdingTypes))

	for _, t := range holdingTypes {
		switch t {
		case types.CharacterEquipType_1HWeapon:
			weightedPicker := util.NewWeightedPicker[table.CreateCharacter1HWeaponRecord, float64](rng)
			for _, record := range v.GetFemale1HWeapon(rng) {
				weightedPicker.Add(record, record.FemaleProb)
			}

			pickedRecord, err := weightedPicker.Pick()
			if err == nil {
				ret[types.CharacterEquipType_1HWeapon] = pickedRecord.Index
			}

		case types.CharacterEquipType_2HWeapon:
			weightedPicker := util.NewWeightedPicker[table.CreateCharacter2HWeaponRecord, float64](rng)
			for _, record := range v.GetFemale2HWeapon(rng) {
				weightedPicker.Add(record, record.FemaleProb)
			}

			pickedRecord, err := weightedPicker.Pick()
			if err == nil {
				ret[types.CharacterEquipType_2HWeapon] = pickedRecord.Index
			}

		case types.CharacterEquipType_SubWeapon:
			weightedPicker := util.NewWeightedPicker[table.CreateCharacterSubWeaponRecord, float64](rng)
			for _, record := range v.GetFemaleSubWeapon(rng) {
				weightedPicker.Add(record, record.FemaleProb)
			}

			pickedRecord, err := weightedPicker.Pick()
			if err == nil {
				ret[types.CharacterEquipType(pickedRecord.SubType)] = pickedRecord.Index
			}
		}
	}

	return ret, len(ret) > 0
}

// 남성 캐릭터 한손 무기를 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterWeaponTableView) GetMale1HWeapon(rng *rand.Rand) []table.CreateCharacter1HWeaponRecord {
	if v.createCharacter1HWeapon == nil {
		return nil
	}

	records := v.createCharacter1HWeapon.GetAll()
	ret := make([]table.CreateCharacter1HWeaponRecord, 0, len(records))
	for _, record := range records {
		if record.MaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}

// 여성 캐릭터 한손 무기를 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterWeaponTableView) GetFemale1HWeapon(rng *rand.Rand) []table.CreateCharacter1HWeaponRecord {
	if v.createCharacter1HWeapon == nil {
		return nil
	}

	records := v.createCharacter1HWeapon.GetAll()
	ret := make([]table.CreateCharacter1HWeaponRecord, 0, len(records))

	for _, record := range records {
		if record.FemaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}

// 남성 캐릭터 양손 무기를 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterWeaponTableView) GetMale2HWeapon(rng *rand.Rand) []table.CreateCharacter2HWeaponRecord {
	if v.createCharacter2HWeapon == nil {
		return nil
	}

	records := v.createCharacter2HWeapon.GetAll()
	ret := make([]table.CreateCharacter2HWeaponRecord, 0, len(records))

	for _, record := range records {
		if record.MaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}

// 여성 캐릭터 양손 무기를 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterWeaponTableView) GetFemale2HWeapon(rng *rand.Rand) []table.CreateCharacter2HWeaponRecord {
	if v.createCharacter2HWeapon == nil {
		return nil
	}

	records := v.createCharacter2HWeapon.GetAll()
	ret := make([]table.CreateCharacter2HWeaponRecord, 0, len(records))

	for _, record := range records {
		if record.FemaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}

// 남성 캐릭터 보조 무기를 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterWeaponTableView) GetMaleSubWeapon(rng *rand.Rand) []table.CreateCharacterSubWeaponRecord {
	if v.createCharacterSubWeapon == nil {
		return nil
	}

	records := v.createCharacterSubWeapon.GetAll()
	ret := make([]table.CreateCharacterSubWeaponRecord, 0, len(records))
	for _, record := range records {
		if record.MaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}

// 여성 캐릭터 보조 무기를 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterWeaponTableView) GetFemaleSubWeapon(rng *rand.Rand) []table.CreateCharacterSubWeaponRecord {
	if v.createCharacterSubWeapon == nil {
		return nil
	}

	records := v.createCharacterSubWeapon.GetAll()
	ret := make([]table.CreateCharacterSubWeaponRecord, 0, len(records))

	for _, record := range records {
		if record.FemaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}
