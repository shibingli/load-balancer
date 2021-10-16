package balancer

import (
	"reflect"

	"github.com/fufuok/load-balancer/utils"
)

type Balancer interface {
	// Select gets next selected item.
	// key is only used for ConsistentHash
	Select(key ...string) interface{}

	// Name load balancer name.
	Name() string

	// Update reinitialize the balancer items.
	Update(choices []*Choice) bool
}

// Choice to be selected for the load balancer
type Choice struct {
	// e.g. server addr / node / *url.URL
	Item interface{}

	// For WeightedRoundRobin / SmoothWeightedRoundRobin / WeightedRand
	Weight int

	// For SmoothWeightedRoundRobin, optional
	CurrentWeight int
}

// Mode defines the selectable balancer algorithm.
type Mode int

const (
	// WeightedRoundRobin is the default balancer algorithm.
	WeightedRoundRobin Mode = iota
	SmoothWeightedRoundRobin
	WeightedRand
	ConsistentHash
	RoundRobin
	Random
)

// NewChoice create new items with optional weights.
func NewChoice(item interface{}, weight ...int) *Choice {
	w := 1
	if len(weight) > 0 {
		w = weight[0]
	}
	return &Choice{
		Item:   item,
		Weight: w,
	}
}

// NewChoicesMap map to choices []*Choice
func NewChoicesMap(items interface{}) (choices []*Choice) {
	v := reflect.ValueOf(items)
	if v.Kind() != reflect.Map {
		return
	}

	n := v.Len()
	m := v.MapRange()
	choices = make([]*Choice, 0, n)
	for m.Next() {
		choices = append(choices, &Choice{
			Item:   m.Key().Interface(),
			Weight: utils.MustInt(m.Value().Interface()),
		})
	}

	return
}

// NewChoicesSlice slice to choices []*Choice
func NewChoicesSlice(items interface{}) (choices []*Choice) {
	v := reflect.ValueOf(items)
	if v.Kind() != reflect.Slice {
		return
	}

	n := v.Len()
	choices = make([]*Choice, 0, n)
	for i := 0; i < n; i++ {
		choices = append(choices, &Choice{
			Item:   v.Index(i).Interface(),
			Weight: 1,
		})
	}

	return
}

// New create a balancer with or without items.
func New(b Mode, choices []*Choice) Balancer {
	switch b {
	case SmoothWeightedRoundRobin:
		return NewSmoothWeightedRoundRobin(choices...)
	case RoundRobin:
		return NewRoundRobin(choices...)
	case WeightedRand:
		return NewWeightedRand(choices...)
	case ConsistentHash:
		return NewConsistentHash(choices...)
	case Random:
		return NewRandom(choices...)
	default:
		return NewWeightedRoundRobin(choices...)
	}
}

// Discard items with a weight less than 1
func cleanWeight(choices []*Choice) (items []*Choice, n int) {
	items = make([]*Choice, 0, len(choices))
	for i := range choices {
		if choices[i].Weight <= 0 {
			continue
		}
		items = append(items, choices[i])
		n++
	}
	return
}
