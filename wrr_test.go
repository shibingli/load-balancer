package balancer

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestWeightedRoundRobin(t *testing.T) {
	lb := NewWeightedRoundRobin()
	item := lb.Select()
	if item != nil {
		t.Fatalf("wrr expected nil, actual %s", item)
	}

	lb = NewWeightedRoundRobin(
		&Choice{Item: "A", Weight: 0},
		&Choice{Item: "B", Weight: 1},
		&Choice{Item: "C", Weight: 7},
		&Choice{Item: "D", Weight: 2},
	)
	count := make(map[string]int)
	for i := 0; i < 1000; i++ {
		item := lb.Select()
		count[item.(string)]++
	}
	if count["A"] != 0 || count["B"] != 100 || count["C"] != 700 || count["D"] != 200 {
		t.Fatal("wrr wrong")
	}

	nodes := []*Choice{
		{Item: "X", Weight: 0},
		{Item: "Y", Weight: 1},
	}
	ok := lb.Update(nodes)
	if ok != true {
		t.Fatal("wrr update wrong")
	}
	item = lb.Select()
	if item != "Y" {
		t.Fatal("wrr update wrong")
	}
}

func TestWeightedRoundRobin_C(t *testing.T) {
	var (
		a, b, c, d int64
		wg         sync.WaitGroup
		mu         sync.Mutex
	)
	nodes := []*Choice{
		{Item: "A", Weight: 5},
		{Item: "B", Weight: 1},
		{Item: "C", Weight: 4},
		{Item: "D", Weight: 0},
	}
	lb := NewWeightedRoundRobin(nodes...)

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 2000; j++ {
				mu.Lock()
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
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if atomic.LoadInt64(&a) != 500000 {
		t.Fatal("wrr wrong: a", atomic.LoadInt64(&a))
	}
	if atomic.LoadInt64(&b) != 100000 {
		t.Fatal("wrr wrong: b", atomic.LoadInt64(&b))
	}
	if atomic.LoadInt64(&c) != 400000 {
		t.Fatal("wrr wrong: c", atomic.LoadInt64(&c))
	}
	if atomic.LoadInt64(&d) != 0 {
		t.Fatal("wrr wrong: d", atomic.LoadInt64(&d))
	}
}
