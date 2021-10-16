package balancer

import (
	"github.com/fufuok/load-balancer/utils"
)

// Weighted Round-Robin Scheduling
// Ref: http://kb.linuxvirtualserver.org/wiki/Weighted_Round-Robin_Scheduling
type wrr struct {
	items []*Choice
	i     int
	n     int
	cw    int
	gcd   int
	max   int
}

func NewWeightedRoundRobin(choices ...*Choice) (lb *wrr) {
	lb = &wrr{}
	lb.Update(choices)
	return
}

func (b *wrr) Select(_ ...string) (item interface{}) {
	switch b.n {
	case 0:
		item = nil
	case 1:
		item = b.items[0].Item
	default:
		item = b.chooseNext().Item
	}
	return
}

func (b *wrr) chooseNext() *Choice {
	for {
		b.i = (b.i + 1) % b.n
		if b.i == 0 {
			b.cw = b.cw - b.gcd
			if b.cw <= 0 {
				b.cw = b.max
				if b.cw == 0 {
					return nil
				}
			}
		}

		if b.items[b.i].Weight >= b.cw {
			return b.items[b.i]
		}
	}
}

func (b *wrr) Name() string {
	return "WeightedRoundRobin"
}

func (b *wrr) Update(choices []*Choice) bool {
	b.items, b.n = cleanWeight(choices)
	b.i = -1
	b.cw = 0
	b.gcd = 0
	b.max = 0

	for i := range b.items {
		b.addSettings(choices[i].Weight)
	}

	return b.n > 0
}

func (b *wrr) addSettings(weight int) {
	if weight > 0 {
		if b.gcd == 0 {
			b.i = -1
			b.cw = 0
			b.gcd = weight
			b.max = weight
		} else {
			b.gcd = utils.GCD(b.gcd, weight)
			if b.max < weight {
				b.max = weight
			}
		}
	}
}
