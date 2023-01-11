package balancer

import (
	"github.com/shibingli/load-balancer/utils"
)

// Random
type random struct {
	items []*Choice
	count uint32
}

func NewRandom(choices ...*Choice) (lb *random) {
	lb = &random{}
	lb.Update(choices)
	return
}

func (b *random) Select(_ ...string) (item interface{}) {
	switch b.count {
	case 0:
		item = nil
	case 1:
		item = b.items[0].Item
	default:
		item = b.items[utils.FastRandn(b.count)].Item
	}
	return
}

func (b *random) Name() string {
	return "Random"
}

func (b *random) Update(choices []*Choice) bool {
	b.items = choices
	b.count = uint32(len(choices))
	return true
}
