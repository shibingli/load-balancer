package balancer

// RoundRobin
type rr struct {
	items   []*Choice
	count   int
	current int
}

func NewRoundRobin(choices ...*Choice) (lb *rr) {
	lb = &rr{}
	lb.Update(choices)
	return
}

func (b *rr) Select(_ ...string) (item interface{}) {
	switch b.count {
	case 0:
		item = nil
	case 1:
		item = b.items[0].Item
	default:
		item = b.items[b.current].Item
		b.current = (b.current + 1) % b.count
	}
	return
}

func (b *rr) Name() string {
	return "RoundRobin"
}

func (b *rr) Update(choices []*Choice) bool {
	b.items = choices
	b.count = len(choices)
	b.current = 0
	return true
}
