package balancer

import (
	"sort"

	"github.com/fufuok/load-balancer/utils"
)

// WeightedRand
type wr struct {
	items   []*Choice
	weights []int
	count   int
	max     uint32
}

func NewWeightedRand(choices ...*Choice) (lb *wr) {
	lb = &wr{}
	lb.Update(choices)
	return
}

func (b *wr) Select(_ ...string) (item interface{}) {
	switch b.count {
	case 0:
		item = nil
	case 1:
		item = b.items[0].Item
	default:
		r := utils.FastRandn(b.max) + 1
		i := utils.SearchInts(b.weights, int(r))
		item = b.items[i].Item
	}
	return
}

func (b *wr) Name() string {
	return "WeightedRand"
}

func (b *wr) Update(choices []*Choice) bool {
	b.items, b.count = cleanWeight(choices)
	sort.Slice(b.items, func(i, j int) bool {
		return b.items[i].Weight < b.items[j].Weight
	})

	max := 0
	weights := make([]int, b.count)
	for i := range b.items {
		max += b.items[i].Weight
		weights[i] = max
	}

	b.weights = weights
	b.max = uint32(max)

	return true
}
