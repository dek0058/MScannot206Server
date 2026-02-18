package view

import (
	"MScannot206/shared/table"
	"MScannot206/shared/types"
	"MScannot206/shared/util"
	"math/rand/v2"
)

// 캐릭터 생성시 장갑을 획득하기 위한 뷰입니다
type CreateCharacterGloveTableView struct {
	createCharacter      *table.CreateCharacterTable
	createCharacterGlove *table.CreateCharacterGloveTable
}

// 캐릭터 장갑을 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterGloveTableView) isHolding(rng *rand.Rand) bool {
	if v.createCharacter == nil || v.createCharacterGlove == nil {
		return false
	}

	ccRecord, ok := v.createCharacter.Get(string(types.CharacterEquipType_Glove))
	if !ok {
		return false
	}

	singleProb := util.NewSingleProbabilisticData[any, float64](1.0, rng)
	singleProb.Set(nil, ccRecord.HoldingProb)

	_, ok = singleProb.Pick()
	return ok
}

// 남성 캐릭터 장갑을 확률적으로 획득합니다
func (v CreateCharacterGloveTableView) GetMale(rng *rand.Rand) (string, bool) {
	if !v.isHolding(rng) {
		return "", false
	}

	weightedPicker := util.NewWeightedPicker[table.CreateCharacterGloveRecord, float64](rng)
	for _, record := range v.GetMaleRecords() {
		weightedPicker.Add(record, record.MaleProb)
	}

	pickedRecord, err := weightedPicker.Pick()
	if err != nil {
		return "", false
	}

	return pickedRecord.Index, true
}

// 여성 캐릭터 장갑을 확률적으로 획득합니다
func (v CreateCharacterGloveTableView) GetFemale(rng *rand.Rand) (string, bool) {
	if !v.isHolding(rng) {
		return "", false
	}

	weightedPicker := util.NewWeightedPicker[table.CreateCharacterGloveRecord, float64](rng)
	for _, record := range v.GetFemaleRecords() {
		weightedPicker.Add(record, record.FemaleProb)
	}

	pickedRecord, err := weightedPicker.Pick()
	if err != nil {
		return "", false
	}

	return pickedRecord.Index, true
}

// 남성 캐릭터 장갑 레코드를 모두 가져옵니다
func (v CreateCharacterGloveTableView) GetMaleRecords() []table.CreateCharacterGloveRecord {
	if v.createCharacterGlove == nil {
		return []table.CreateCharacterGloveRecord{}
	}

	records := v.createCharacterGlove.GetAll()
	ret := make([]table.CreateCharacterGloveRecord, 0, len(records))
	for _, record := range records {
		if record.MaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}

// 여성 캐릭터 장갑 레코드를 모두 가져옵니다
func (v CreateCharacterGloveTableView) GetFemaleRecords() []table.CreateCharacterGloveRecord {
	if v.createCharacterGlove == nil {
		return []table.CreateCharacterGloveRecord{}
	}

	records := v.createCharacterGlove.GetAll()
	ret := make([]table.CreateCharacterGloveRecord, 0, len(records))
	for _, record := range records {
		if record.FemaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}
