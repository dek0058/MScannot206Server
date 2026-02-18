package view

import (
	"MScannot206/shared/table"
	"MScannot206/shared/types"
	"MScannot206/shared/util"
	"math/rand/v2"
)

// 캐릭터 생성시 모자를 획득하기 위한 뷰입니다
type CreateCharacterCapTableView struct {
	createCharacter    *table.CreateCharacterTable
	createCharacterCap *table.CreateCharacterCapTable
}

// 캐릭터 모자를 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterCapTableView) isHolding(rng *rand.Rand) bool {
	if v.createCharacter == nil || v.createCharacterCap == nil {
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
func (v CreateCharacterCapTableView) GetMale(rng *rand.Rand) (string, bool) {
	if !v.isHolding(rng) {
		return "", false
	}

	weightedPicker := util.NewWeightedPicker[table.CreateCharacterCapRecord, float64](rng)
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
func (v CreateCharacterCapTableView) GetFemale(rng *rand.Rand) (string, bool) {
	if !v.isHolding(rng) {
		return "", false
	}

	weightedPicker := util.NewWeightedPicker[table.CreateCharacterCapRecord, float64](rng)
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
func (v CreateCharacterCapTableView) GetMaleRecords() []table.CreateCharacterCapRecord {
	if v.createCharacterCap == nil {
		return []table.CreateCharacterCapRecord{}
	}

	records := v.createCharacterCap.GetAll()
	ret := make([]table.CreateCharacterCapRecord, 0, len(records))
	for _, record := range records {
		if record.MaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}

// 여성 캐릭터 모자 레코드를 모두 가져옵니다
func (v CreateCharacterCapTableView) GetFemaleRecords() []table.CreateCharacterCapRecord {
	if v.createCharacterCap == nil {
		return []table.CreateCharacterCapRecord{}
	}

	records := v.createCharacterCap.GetAll()
	ret := make([]table.CreateCharacterCapRecord, 0, len(records))
	for _, record := range records {
		if record.FemaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}
