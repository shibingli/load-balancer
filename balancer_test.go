package balancer

import (
	"testing"
)

func TestBalancer(t *testing.T) {
	lb := New(WeightedRoundRobin, nil)
	if lb.Name() != "WeightedRoundRobin" {
		t.Fatal("balancer.New wrong")
	}

	lb = New(SmoothWeightedRoundRobin, nil)
	if lb.Name() != "SmoothWeightedRoundRobin" {
		t.Fatal("balancer.New wrong")
	}

	lb = New(WeightedRand, nil)
	if lb.Name() != "WeightedRand" {
		t.Fatal("balancer.New wrong")
	}

	wNodes := map[interface{}]int{
		"X": 0,
		"Y": 1,
	}
	choices := NewChoicesMap(wNodes)
	lb.Update(choices)
	best := lb.Select()
	if best != "Y" {
		t.Fatal("balancer select wrong")
	}

	lb = New(ConsistentHash, nil)
	if lb.Name() != "ConsistentHash" {
		t.Fatal("balancer.New wrong")
	}

	lb = New(RoundRobin, nil)
	if lb.Name() != "RoundRobin" {
		t.Fatal("balancer.New wrong")
	}

	lb = New(Random, nil)
	if lb.Name() != "Random" {
		t.Fatal("balancer.New wrong")
	}

	lb.Update([]*Choice{
		NewChoice("A"),
	})
	best = lb.Select()
	if best != "A" {
		t.Fatal("balancer select wrong")
	}

	nodes := []string{"B", "C"}
	choices = NewChoicesSlice(nodes)
	lb.Update(choices)
	best = lb.Select()
	if best != "B" && best != "C" {
		t.Fatal("balancer select wrong")
	}
}
