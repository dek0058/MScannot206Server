package view

import (
	"MScannot206/shared/table"
	"MScannot206/shared/types"
	"MScannot206/shared/util"
	"math/rand/v2"
)

// 캐릭터 생성시 머리를 획득하기 위한 뷰입니다
type CreateCharacterHairTableView struct {
	CreateCharacter     *table.CreateCharacterTable
	CreateCharacterHair *table.CreateCharacterHairTable
}

// 캐릭터 머리를 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterHairTableView) isHolding(rng *rand.Rand) bool {
	if v.CreateCharacter == nil || v.CreateCharacterHair == nil {
		return false
	}

	ccRecord, ok := v.CreateCharacter.Get(string(types.CharacterEquipType_Hair))
	if !ok {
		return false
	}

	singleProb := util.NewSingleProbabilisticData[any](1.0, rng)
	singleProb.Set(nil, ccRecord.HoldingProb)

	_, ok = singleProb.Pick()
	return ok
}

// 남성 캐릭터 머리를 확률적으로 획득합니다
func (v CreateCharacterHairTableView) GetMale(rng *rand.Rand) (string, bool) {
	if !v.isHolding(rng) {
		return "", false
	}

	weightedPicker := util.NewWeightedPicker[table.CreateCharacterHairRecord, float64](rng)
	for _, record := range v.GetMaleRecords() {
		weightedPicker.Add(record, record.MaleProb)
	}

	pickedRecord, err := weightedPicker.Pick()
	if err != nil {
		return "", false
	}

	return pickedRecord.Index, true
}

// 여성 캐릭터 머리를 확률적으로 획득합니다
func (v CreateCharacterHairTableView) GetFemale(rng *rand.Rand) (string, bool) {
	if !v.isHolding(rng) {
		return "", false
	}

	weightedPicker := util.NewWeightedPicker[table.CreateCharacterHairRecord, float64](rng)
	for _, record := range v.GetFemaleRecords() {
		weightedPicker.Add(record, record.FemaleProb)
	}

	pickedRecord, err := weightedPicker.Pick()
	if err != nil {
		return "", false
	}

	return pickedRecord.Index, true
}

// 남성 캐릭터 머리 레코드를 모두 가져옵니다
func (v CreateCharacterHairTableView) GetMaleRecords() []table.CreateCharacterHairRecord {
	if v.CreateCharacterHair == nil {
		return []table.CreateCharacterHairRecord{}
	}

	records := v.CreateCharacterHair.GetAll()
	ret := make([]table.CreateCharacterHairRecord, 0, len(records))
	for _, record := range records {
		if record.MaleProb > 0.0 {
			ret = append(ret, record)
		}
	}

	return ret
}

// 여성 캐릭터 머리 레코드를 모두 가져옵니다
func (v CreateCharacterHairTableView) GetFemaleRecords() []table.CreateCharacterHairRecord {
	if v.CreateCharacterHair == nil {
		return []table.CreateCharacterHairRecord{}
	}

	records := v.CreateCharacterHair.GetAll()
	ret := make([]table.CreateCharacterHairRecord, 0, len(records))
	for _, record := range records {
		if record.FemaleProb > 0.0 {
			ret = append(ret, record)
		}
	}

	return ret
}
