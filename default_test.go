package balancer

import (
	"testing"
)

func TestDefaultBalancer(t *testing.T) {
	item := Select()
	if item != nil {
		t.Fatalf("default balancer expected nil, actual %s", item)
	}

	nodes := []*Choice{
		{Item: "A", Weight: 0},
		{Item: "B", Weight: 1},
		{Item: "C", Weight: 7},
		{Item: "D", Weight: 2},
	}
	Update(nodes)
	count := make(map[string]int)
	for i := 0; i < 1000; i++ {
		item := Select()
		count[item.(string)]++
	}
	if count["A"] != 0 || count["B"] != 100 || count["C"] != 700 || count["D"] != 200 {
		t.Fatal("default balancer wrong")
	}

	nodes = []*Choice{
		{Item: "X", Weight: 0},
		{Item: "Y", Weight: 1},
	}
	ok := Update(nodes)
	if ok != true {
		t.Fatal("default balancer update wrong")
	}
	item = Select()
	if item != "Y" {
		t.Fatal("default balancer update wrong")
	}

	if Name() != "WeightedRoundRobin" {
		t.Fatal("default balancer name wrong")
	}
}
