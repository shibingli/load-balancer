package balancer

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRandom(t *testing.T) {
	lb := NewRandom()
	item := lb.Select()
	if item != nil {
		t.Fatalf("r expected nil, actual %s", item)
	}

	lb = NewRandom(
		&Choice{Item: "A"},
		&Choice{Item: "B"},
		&Choice{Item: "C"},
		&Choice{Item: "D"},
	)
	item = lb.Select()
	count := make(map[string]int)
	for i := 0; i < 2000; i++ {
		item := lb.Select()
		count[item.(string)]++
	}
	if count["A"] <= 300 || count["B"] <= 300 || count["C"] <= 300 || count["D"] <= 300 {
		t.Fatal("r wrong")
	}
	if count["A"]+count["B"]+count["C"]+count["D"] != 2000 {
		t.Fatal("r wrong")
	}

	nodes := []*Choice{
		{Item: "X"},
		{Item: "Y"},
	}
	ok := lb.Update(nodes)
	if ok != true {
		t.Fatal("r update wrong")
	}
	item = lb.Select()
	if item != "X" && item != "Y" {
		t.Fatal("r update wrong")
	}
}

func TestRandom_C(t *testing.T) {
	var (
		a, b, c, d int64
	)
	nodes := []*Choice{
		{Item: "A"},
		{Item: "B"},
		{Item: "C"},
		{Item: "D"},
	}
	lb := NewRandom(nodes...)

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
		t.Fatal("r wrong: a")
	}
	if atomic.LoadInt64(&b) <= 200000 {
		t.Fatal("r wrong: b")
	}
	if atomic.LoadInt64(&c) <= 200000 {
		t.Fatal("r wrong: c")
	}
	if atomic.LoadInt64(&d) <= 200000 {
		t.Fatal("r wrong: d")
	}
	if atomic.LoadInt64(&a)+atomic.LoadInt64(&b)+atomic.LoadInt64(&c)+atomic.LoadInt64(&d) != 1000000 {
		t.Fatal("r wrong: sum")
	}
}
