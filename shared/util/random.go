package util

import (
	"errors"
	"math/rand/v2"
)

// Number는 가중치로 사용할 수 있는 숫자 타입 제약 조건입니다.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Item은 아이템과 가중치를 저장하는 내부 구조체입니다.
type Item[T any, W Number] struct {
	Data   T
	Weight W
}

// 1. WeightedPicker: Data를 넣으면 모든 Data의 weight를 더하고 각 weight의 값 중 하나를 선택하는 랜덤
type WeightedPicker[T any, W Number] struct {
	items       []Item[T, W]
	totalWeight W
	rng         *rand.Rand
}

func NewWeightedPicker[T any, W Number](rng *rand.Rand) *WeightedPicker[T, W] {
	return &WeightedPicker[T, W]{
		items: make([]Item[T, W], 0),
		rng:   rng,
	}
}

func (w *WeightedPicker[T, W]) Add(data T, weight W) {
	if weight <= 0 {
		return
	}
	w.items = append(w.items, Item[T, W]{Data: data, Weight: weight})
	w.totalWeight += weight
}

func (w *WeightedPicker[T, W]) Pick() (T, error) {
	var empty T
	if w.totalWeight == 0 {
		return empty, errors.New("total weight is 0")
	}

	r := w.randomValue(w.totalWeight)
	var current W
	for _, item := range w.items {
		current += item.Weight
		if r < float64(current) {
			return item.Data, nil
		}
	}
	return empty, errors.New("unexpected error in Pick")
}

func (w *WeightedPicker[T, W]) randomValue(max W) float64 {
	// W가 정수형인지 실수형인지에 따라 다르게 처리하지 않고,
	// 단순히 float64로 변환하여 0 ~ max(exclusive) 사이의 값을 구함
	// 정수형의 경우 정확한 인덱싱이 필요하다면 별도 로직이 필요할 수 있으나,
	// 가중치 뽑기(룰렛) 로직에서는 float64 변환 후 비교해도 무방함.
	return w.rng.Float64() * float64(max)
}

// 2. IndependentPicker: Data를 넣으면 Data가 가진 weight를 각각 랜덤을 돌려서 나온 Data를 모두 리턴 하는 랜덤
// maxWeight는 확률의 기준값입니다 (예: 100, 10000, 1.0 등).
type IndependentPicker[T any, W Number] struct {
	items     []Item[T, W]
	maxWeight W
	rng       *rand.Rand
}

func NewIndependentPicker[T any, W Number](maxWeight W, rng *rand.Rand) *IndependentPicker[T, W] {
	return &IndependentPicker[T, W]{
		items:     make([]Item[T, W], 0),
		maxWeight: maxWeight,
		rng:       rng,
	}
}

func (i *IndependentPicker[T, W]) Add(data T, weight W) {
	if weight <= 0 {
		return
	}
	i.items = append(i.items, Item[T, W]{Data: data, Weight: weight})
}

func (i *IndependentPicker[T, W]) Pick() []T {
	result := make([]T, 0)
	for _, item := range i.items {
		// 0 <= rand < maxWeight
		if i.rng.Float64()*float64(i.maxWeight) < float64(item.Weight) {
			result = append(result, item.Data)
		}
	}
	return result
}

// 3. DeckPicker: Data를 넣으면 모든 Data의 weight를 더하고 각 weight의 값 중 하나를 선택하고 선택된 Data는 삭제하는 랜덤
type DeckPicker[T any, W Number] struct {
	items       []Item[T, W]
	totalWeight W
	rng         *rand.Rand
}

func NewDeckPicker[T any, W Number](rng *rand.Rand) *DeckPicker[T, W] {
	return &DeckPicker[T, W]{
		items: make([]Item[T, W], 0),
		rng:   rng,
	}
}

func (d *DeckPicker[T, W]) Add(data T, weight W) {
	if weight <= 0 {
		return
	}
	d.items = append(d.items, Item[T, W]{Data: data, Weight: weight})
	d.totalWeight += weight
}

func (d *DeckPicker[T, W]) Pick() (T, bool) {
	var empty T
	if d.totalWeight == 0 {
		return empty, false
	}

	r := d.rng.Float64() * float64(d.totalWeight)
	var current W
	selectedIndex := -1
	var selectedData T

	// 선택된 아이템 찾기
	for i, item := range d.items {
		current += item.Weight
		if r < float64(current) {
			selectedData = item.Data
			selectedIndex = i
			break
		}
	}

	if selectedIndex == -1 {
		return empty, false
	}

	// 선택된 아이템 제거 및 전체 가중치 업데이트
	d.totalWeight -= d.items[selectedIndex].Weight
	d.items = append(d.items[:selectedIndex], d.items[selectedIndex+1:]...)

	return selectedData, true
}

// 4. SingleProbabilisticData: 데이터 하나와 가중치(확률), 그리고 최대값을 설정하여 확률적으로 아이템을 획득하는 랜덤
// 예: 가중치 0.6, 최대값 1.0 설정 시 60% 확률로 아이템 획득
type SingleProbabilisticData[T any, W Number] struct {
	data      T
	weight    W
	maxWeight W
	rng       *rand.Rand
	hasData   bool
}

func NewSingleProbabilisticData[T any, W Number](maxWeight W, rng *rand.Rand) *SingleProbabilisticData[T, W] {
	return &SingleProbabilisticData[T, W]{
		maxWeight: maxWeight,
		rng:       rng,
		hasData:   false,
	}
}

func (s *SingleProbabilisticData[T, W]) Set(data T, weight W) {
	s.data = data
	s.weight = weight
	s.hasData = true
}

func (s *SingleProbabilisticData[T, W]) Pick() (T, bool) {
	var empty T
	if !s.hasData || s.maxWeight <= 0 {
		return empty, false
	}

	if s.rng.Float64()*float64(s.maxWeight) < float64(s.weight) {
		return s.data, true
	}

	return empty, false
}
