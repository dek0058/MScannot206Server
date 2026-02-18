package view

import (
	"MScannot206/shared/table"
	"MScannot206/shared/types"
	"MScannot206/shared/util"
	"math/rand/v2"
)

type CreateCharacterCapeTableView struct {
	createCharacter     *table.CreateCharacterTable
	createCharacterCape *table.CreateCharacterCapeTable
}

// 캐릭터 모자를 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterCapeTableView) isHolding(rng *rand.Rand) bool {
	if v.createCharacter == nil || v.createCharacterCape == nil {
		return false
	}

	ccRecord, ok := v.createCharacter.Get(string(types.CharacterEquipType_Cap))
	if !ok {
		return false
	}

	singleProb := util.NewSingleProbabilisticData[any](1.0, rng)
	singleProb.Set(nil, ccRecord.HoldingProb)

	_, ok = singleProb.Pick()
	return ok
}

// 남성 캐릭터 모자를 확률적으로 획득합니다
func (v CreateCharacterCapeTableView) GetMale(rng *rand.Rand) (string, bool) {
	if !v.isHolding(rng) {
		return "", false
	}

	weightedPicker := util.NewWeightedPicker[table.CreateCharacterCapeRecord, float64](rng)
	for _, record := range v.GetMaleRecords() {
		weightedPicker.Add(record, record.MaleProb)
	}

	pickedRecord, err := weightedPicker.Pick()
	if err != nil {
		return "", false
	}

	return pickedRecord.Index, true
}

// 여성 캐릭터 모자를 확률적으로 획득합니다
func (v CreateCharacterCapeTableView) GetFemale(rng *rand.Rand) (string, bool) {
	if !v.isHolding(rng) {
		return "", false
	}

	weightedPicker := util.NewWeightedPicker[table.CreateCharacterCapeRecord, float64](rng)
	for _, record := range v.GetFemaleRecords() {
		weightedPicker.Add(record, record.FemaleProb)
	}

	pickedRecord, err := weightedPicker.Pick()
	if err != nil {
		return "", false
	}

	return pickedRecord.Index, true
}

// 남성 캐릭터 모자 레코드를 모두 가져옵니다
func (v CreateCharacterCapeTableView) GetMaleRecords() []table.CreateCharacterCapeRecord {
	if v.createCharacterCape == nil {
		return []table.CreateCharacterCapeRecord{}
	}

	records := v.createCharacterCape.GetAll()
	ret := make([]table.CreateCharacterCapeRecord, 0, len(records))
	for _, record := range records {
		if record.MaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}

// 여성 캐릭터 모자 레코드를 모두 가져옵니다
func (v CreateCharacterCapeTableView) GetFemaleRecords() []table.CreateCharacterCapeRecord {
	if v.createCharacterCape == nil {
		return []table.CreateCharacterCapeRecord{}
	}

	records := v.createCharacterCape.GetAll()
	ret := make([]table.CreateCharacterCapeRecord, 0, len(records))

	for _, record := range records {
		if record.FemaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}
