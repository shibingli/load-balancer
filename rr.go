package balancer

import (
	"sync/atomic"
)

// RoundRobin
type rr struct {
	items   []*Choice
	count   uint32
	current uint32
}

func NewRoundRobin(choices ...*Choice) (lb *rr) {
	lb = &rr{}
	lb.Update(choices)
	return
}

func (b *rr) Select(_ ...string) (item interface{}) {
	n := atomic.LoadUint32(&b.count)
	switch n {
	case 0:
		item = nil
	case 1:
		item = b.items[0].Item
	default:
		m := atomic.LoadUint32(&b.current)
		item = b.items[m].Item
		m = (m + 1) % n
		atomic.StoreUint32(&b.current, m)
	}
	return
}

func (b *rr) Name() string {
	return "RoundRobin"
}

func (b *rr) Update(choices []*Choice) bool {
	b.items = choices
	b.count = uint32(len(choices))
	b.current = 0
	return true
}
