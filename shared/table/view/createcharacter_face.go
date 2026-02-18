package view

import (
	"MScannot206/shared/table"
	"MScannot206/shared/types"
	"MScannot206/shared/util"
	"math/rand/v2"
)

// 캐릭터 생성시 얼굴을 획득하기 위한 뷰입니다
type CreateCharacterFaceTableView struct {
	CreateCharacter     *table.CreateCharacterTable
	CreateCharacterFace *table.CreateCharacterFaceTable
}

// 캐릭터 얼굴을 획득 할 수 있는지 확률을 통해 판단합니다
func (v CreateCharacterFaceTableView) isHolding(rng *rand.Rand) bool {
	if v.CreateCharacter == nil || v.CreateCharacterFace == nil {
		return false
	}

	ccRecord, ok := v.CreateCharacter.Get(string(types.CharacterEquipType_Face))
	if !ok {
		return false
	}

	singleProb := util.NewSingleProbabilisticData[any](1.0, rng)
	singleProb.Set(nil, ccRecord.HoldingProb)

	_, ok = singleProb.Pick()
	return ok
}

// 남성 캐릭터 얼굴을 확률적으로 획득합니다
func (v CreateCharacterFaceTableView) GetMale(rng *rand.Rand) (string, bool) {
	if !v.isHolding(rng) {
		return "", false
	}

	weightedPicker := util.NewWeightedPicker[table.CreateCharacterFaceRecord, float64](rng)
	for _, record := range v.GetMaleRecords() {
		weightedPicker.Add(record, record.MaleProb)
	}

	pickedRecord, err := weightedPicker.Pick()
	if err != nil {
		return "", false
	}

	return pickedRecord.Index, true
}

// 여성 캐릭터 얼굴을 확률적으로 획득합니다
func (v CreateCharacterFaceTableView) GetFemale(rng *rand.Rand) (string, bool) {
	if !v.isHolding(rng) {
		return "", false
	}

	weightedPicker := util.NewWeightedPicker[table.CreateCharacterFaceRecord, float64](rng)
	for _, record := range v.GetFemaleRecords() {
		weightedPicker.Add(record, record.FemaleProb)
	}

	pickedRecord, err := weightedPicker.Pick()
	if err != nil {
		return "", false
	}

	return pickedRecord.Index, true
}

// 남성 캐릭터 얼굴 레코드를 모두 가져옵니다
func (v CreateCharacterFaceTableView) GetMaleRecords() []table.CreateCharacterFaceRecord {
	if v.CreateCharacterFace == nil {
		return []table.CreateCharacterFaceRecord{}
	}

	records := v.CreateCharacterFace.GetAll()
	ret := make([]table.CreateCharacterFaceRecord, 0, len(records))
	for _, record := range records {
		if record.MaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}

// 여성 캐릭터 얼굴 레코드를 모두 가져옵니다
func (v CreateCharacterFaceTableView) GetFemaleRecords() []table.CreateCharacterFaceRecord {
	if v.CreateCharacterFace == nil {
		return []table.CreateCharacterFaceRecord{}
	}

	records := v.CreateCharacterFace.GetAll()
	ret := make([]table.CreateCharacterFaceRecord, 0, len(records))
	for _, record := range records {
		if record.FemaleProb > 0 {
			ret = append(ret, record)
		}
	}

	return ret
}
