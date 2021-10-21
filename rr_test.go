package balancer

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRoundRobin(t *testing.T) {
	lb := NewRoundRobin()
	item := lb.Select()
	if item != nil {
		t.Fatalf("rr expected nil, actual %s", item)
	}

	lb = NewRoundRobin(
		&Choice{Item: "A"},
		&Choice{Item: "B"},
		&Choice{Item: "C"},
		&Choice{Item: "D"},
	)
	item = lb.Select()
	if item != "A" {
		t.Fatalf("rr expected A, actual %s", item)
	}
	item = lb.Select()
	if item != "B" {
		t.Fatalf("rr expected B, actual %s", item)
	}
	item = lb.Select()
	if item != "C" {
		t.Fatalf("rr expected C, actual %s", item)
	}
	item = lb.Select()
	if item != "D" {
		t.Fatalf("rr expected D, actual %s", item)
	}
	item = lb.Select()
	if item != "A" {
		t.Fatalf("rr expected A, actual %s", item)
	}

	count := make(map[string]int)
	for i := 0; i < 2000; i++ {
		item := lb.Select()
		count[item.(string)]++
	}
	if count["A"] != 500 || count["B"] != 500 || count["C"] != 500 || count["D"] != 500 {
		t.Fatal("rr wrong")
	}

	nodes := []*Choice{
		{Item: "X"},
		{Item: "Y"},
	}
	ok := lb.Update(nodes)
	if ok != true {
		t.Fatal("rr update wrong")
	}
	item = lb.Select()
	if item != "X" {
		t.Fatal("rr update wrong")
	}
}

func TestRoundRobin_C(t *testing.T) {
	var (
		a, b, c, d int64
	)
	nodes := []*Choice{
		{Item: "A"},
		{Item: "B"},
		{Item: "C"},
		{Item: "D"},
	}
	lb := NewRoundRobin(nodes...)

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

	t.Log(atomic.LoadInt64(&a), atomic.LoadInt64(&b), atomic.LoadInt64(&c), atomic.LoadInt64(&d))

	if atomic.LoadInt64(&a) <= 200000 {
		t.Fatal("rr wrong: a")
	}
	if atomic.LoadInt64(&b) <= 200000 {
		t.Fatal("rr wrong: b")
	}
	if atomic.LoadInt64(&c) <= 200000 {
		t.Fatal("rr wrong: c")
	}
	if atomic.LoadInt64(&d) <= 200000 {
		t.Fatal("rr wrong: d")
	}
	if atomic.LoadInt64(&a)+atomic.LoadInt64(&b)+atomic.LoadInt64(&c)+atomic.LoadInt64(&d) != 1000000 {
		t.Fatal("rr wrong: sum")
	}
}
