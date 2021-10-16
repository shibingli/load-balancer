package balancer

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestWeightedRand(t *testing.T) {
	lb := NewWeightedRand()
	item := lb.Select()
	if item != nil {
		t.Fatalf("wr expected nil, actual %s", item)
	}

	lb = NewWeightedRand(
		&Choice{Item: "A", Weight: 0},
		&Choice{Item: "B", Weight: 1},
		&Choice{Item: "C", Weight: 7},
		&Choice{Item: "D", Weight: 2},
	)
	count := make(map[string]int)
	for i := 0; i < 2000; i++ {
		item := lb.Select()
		count[item.(string)]++
	}
	if count["A"] != 0 || count["B"] <= 150 || count["C"] <= 750 || count["D"] <= 250 {
		t.Fatal("wr wrong")
	}
	if count["A"]+count["B"]+count["C"]+count["D"] != 2000 {
		t.Fatal("wr wrong")
	}

	nodes := []*Choice{
		{Item: "X", Weight: 0},
		{Item: "Y", Weight: 1},
	}
	ok := lb.Update(nodes)
	if ok != true {
		t.Fatal("wr update wrong")
	}
	item = lb.Select()
	if item != "Y" {
		t.Fatal("wr update wrong")
	}
}

func TestWeightedRand_C(t *testing.T) {
	var (
		a, b, c, d int64
	)
	nodes := []*Choice{
		{Item: "A", Weight: 5},
		{Item: "B", Weight: 1},
		{Item: "C", Weight: 4},
		{Item: "D", Weight: 0},
	}
	lb := NewWeightedRand(nodes...)

	var wg sync.WaitGroup
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 2000; j++ {
				switch lb.Select() {
				case "A":
					atomic.AddInt64(&a, 1)
				case "B":
					atomic.AddInt64(&b, 1)
				case "C":
					atomic.AddInt64(&c, 1)
				case "D":
					atomic.AddInt64(&d, 1)
				}
			}
		}()
	}
	wg.Wait()

	if atomic.LoadInt64(&a) <= 400000 {
		t.Fatal("wr wrong: a")
	}
	if atomic.LoadInt64(&b) <= 40000 {
		t.Fatal("wr wrong: b")
	}
	if atomic.LoadInt64(&c) <= 300000 {
		t.Fatal("wr wrong: c")
	}
	if atomic.LoadInt64(&d) != 0 {
		t.Fatal("wr wrong: d")
	}
	if atomic.LoadInt64(&a)+atomic.LoadInt64(&b)+atomic.LoadInt64(&c) != 1000000 {
		t.Fatal("wr wrong: sum")
	}
}
