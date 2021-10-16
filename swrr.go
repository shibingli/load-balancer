package balancer

// Smooth weighted round-robin balancing
// Ref: https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35
type swrr struct {
	items []*Choice
	count int
}

func NewSmoothWeightedRoundRobin(choices ...*Choice) (lb *swrr) {
	lb = &swrr{}
	lb.Update(choices)
	return
}

func (b *swrr) Select(_ ...string) (item interface{}) {
	switch b.count {
	case 0:
		item = nil
	case 1:
		item = b.items[0].Item
	default:
		item = b.chooseNext().Item
	}
	return
}

func (b *swrr) chooseNext() (choice *Choice) {
	total := 0
	for i := range b.items {
		c := b.items[i]
		if c == nil {
			return nil
		}

		total += c.Weight
		c.CurrentWeight += c.Weight

		if choice == nil || c.CurrentWeight > choice.CurrentWeight {
			choice = c
		}
	}

	if choice == nil {
		return nil
	}

	choice.CurrentWeight -= total

	return choice
}

func (b *swrr) Name() string {
	return "SmoothWeightedRoundRobin"
}

func (b *swrr) Update(choices []*Choice) bool {
	b.items, b.count = cleanWeight(choices)
	return b.count > 0
}
