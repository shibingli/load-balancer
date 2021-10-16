package balancer

import (
	"github.com/fufuok/load-balancer/internal/doublejump"
	"github.com/fufuok/load-balancer/utils"
)

// JumpConsistentHash
type consistentHash struct {
	count int
	h     *doublejump.Hash
}

func NewConsistentHash(choices ...*Choice) (lb *consistentHash) {
	lb = &consistentHash{}
	lb.Update(choices)
	return
}

func (b *consistentHash) Select(key ...string) (item interface{}) {
	if b.count == 0 {
		return
	}
	hash := utils.HashString(key...)
	return b.h.Get(hash).(*Choice).Item
}

func (b *consistentHash) Name() string {
	return "ConsistentHash"
}

func (b *consistentHash) Update(choices []*Choice) bool {
	b.count = len(choices)
	b.h = doublejump.NewHash()
	for i := range choices {
		b.h.Add(choices[i])
	}
	return true
}
