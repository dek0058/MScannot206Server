package view

import (
	"MScannot206/shared/table"
	"MScannot206/shared/types"
	"MScannot206/shared/util"
	"math/rand/v2"
)

// 캐릭터 생성시 귀 모양을 획득하기 위한 뷰입니다
type CreateCharacterEarTableView struct {
	createCharacter    *table.CreateCharacterTable
	createCharacterEar *table.CreateCharacterEarTable
}

// 캐릭터 귀 모양을 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterEarTableView) isHolding(rng *rand.Rand) bool {
	if v.createCharacter == nil || v.createCharacterEar == nil {
		return false
	}

	ccRecord, ok := v.createCharacter.Get(string(types.CharacterEquipType_Ear))
	if !ok {
		return false
	}

	singleProb := util.NewSingleProbabilisticData[any, float64](1.0, rng)
	singleProb.Set(nil, ccRecord.HoldingProb)

	_, ok = singleProb.Pick()
	return ok
}

// 남성 캐릭터 귀 모양을 확률적으로 획득합니다
func (v CreateCharacterEarTableView) GetMale(rng *rand.Rand) (string, bool) {
	if !v.isHolding(rng) {
		return "", false
	}

	weightedPicker := util.NewWeightedPicker[table.CreateCharacterEarRecord, float64](rng)
	for _, record := range v.GetMaleRecords() {
		weightedPicker.Add(record, record.MaleProb)
	}

	pickedRecord, err := weightedPicker.Pick()
	if err != nil {
		return "", false
	}

	return pickedRecord.Index, true
}

// 여성 캐릭터 귀 모양을 확률적으로 획득합니다
func (v CreateCharacterEarTableView) GetFemale(rng *rand.Rand) (string, bool) {
	if !v.isHolding(rng) {
		return "", false
	}

	weightedPicker := util.NewWeightedPicker[table.CreateCharacterEarRecord, float64](rng)
	for _, record := range v.GetFemaleRecords() {
		weightedPicker.Add(record, record.FemaleProb)
	}

	pickedRecord, err := weightedPicker.Pick()
	if err != nil {
		return "", false
	}

	return pickedRecord.Index, true
}

// 남성 캐릭터 귀 모양 레코드를 모두 가져옵니다
func (v CreateCharacterEarTableView) GetMaleRecords() []table.CreateCharacterEarRecord {
	if v.createCharacterEar == nil {
		return []table.CreateCharacterEarRecord{}
	}

	records := v.createCharacterEar.GetAll()
	ret := make([]table.CreateCharacterEarRecord, 0, len(records))
	for _, record := range records {
		if record.MaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}

// 여성 캐릭터 귀 모양 레코드를 모두 가져옵니다
func (v CreateCharacterEarTableView) GetFemaleRecords() []table.CreateCharacterEarRecord {
	if v.createCharacterEar == nil {
		return []table.CreateCharacterEarRecord{}
	}

	records := v.createCharacterEar.GetAll()
	ret := make([]table.CreateCharacterEarRecord, 0, len(records))
	for _, record := range records {
		if record.FemaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}
