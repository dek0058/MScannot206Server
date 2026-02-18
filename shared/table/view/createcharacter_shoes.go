package view

import (
	"MScannot206/shared/table"
	"MScannot206/shared/types"
	"MScannot206/shared/util"
	"math/rand/v2"
)

// 캐릭터 생성시 신발을 획득하기 위한 뷰입니다
type CreateCharacterShoesTableView struct {
	createCharacter      *table.CreateCharacterTable
	createCharacterShoes *table.CreateCharacterShoesTable
}

// 캐릭터 신발을 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterShoesTableView) isHolding(rng *rand.Rand) bool {
	if v.createCharacter == nil || v.createCharacterShoes == nil {
		return false
	}

	ccRecord, ok := v.createCharacter.Get(string(types.CharacterEquipType_Shoes))
	if !ok {
		return false
	}

	singleProb := util.NewSingleProbabilisticData[any, float64](1.0, rng)
	singleProb.Set(nil, ccRecord.HoldingProb)

	_, ok = singleProb.Pick()
	return ok
}

// 남성 캐릭터 신발을 확률적으로 획득합니다
func (v CreateCharacterShoesTableView) GetMale(rng *rand.Rand) (string, bool) {
	if !v.isHolding(rng) {
		return "", false
	}

	weightedPicker := util.NewWeightedPicker[table.CreateCharacterShoesRecord, float64](rng)
	for _, record := range v.GetMaleRecords() {
		weightedPicker.Add(record, record.MaleProb)
	}

	pickedRecord, err := weightedPicker.Pick()
	if err != nil {
		return "", false
	}

	return pickedRecord.Index, true
}

// 여성 캐릭터 신발을 확률적으로 획득합니다
func (v CreateCharacterShoesTableView) GetFemale(rng *rand.Rand) (string, bool) {
	if !v.isHolding(rng) {
		return "", false
	}

	weightedPicker := util.NewWeightedPicker[table.CreateCharacterShoesRecord, float64](rng)
	for _, record := range v.GetFemaleRecords() {
		weightedPicker.Add(record, record.FemaleProb)
	}

	pickedRecord, err := weightedPicker.Pick()
	if err != nil {
		return "", false
	}

	return pickedRecord.Index, true
}

// 남성 캐릭터 신발 레코드를 모두 가져옵니다
func (v CreateCharacterShoesTableView) GetMaleRecords() []table.CreateCharacterShoesRecord {
	if v.createCharacterShoes == nil {
		return []table.CreateCharacterShoesRecord{}
	}

	records := v.createCharacterShoes.GetAll()
	ret := make([]table.CreateCharacterShoesRecord, 0, len(records))
	for _, record := range records {
		if record.MaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}

// 여성 캐릭터 신발 레코드를 모두 가져옵니다
func (v CreateCharacterShoesTableView) GetFemaleRecords() []table.CreateCharacterShoesRecord {
	if v.createCharacterShoes == nil {
		return []table.CreateCharacterShoesRecord{}
	}

	records := v.createCharacterShoes.GetAll()
	ret := make([]table.CreateCharacterShoesRecord, 0, len(records))
	for _, record := range records {
		if record.FemaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}
