# ⚖️ load balancing algorithm library

High-performance general load balancing algorithm library, non-goroutine-safe.

Smooth weighted load balancing algorithm: [NGINX](https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35) and [LVS](http://kb.linuxvirtualserver.org/wiki/Weighted_Round-Robin_Scheduling), Doublejump provides a revamped Google's jump consistent hash.

------

If you want a **goroutine-safe** load balancer, please refer to another library, which supports more APIs: [fufuok/balancer](https://github.com/fufuok/load-balancer)

------

## 🎯 Features

- WeightedRoundRobin
- SmoothWeightedRoundRobin
- WeightedRand
- ConsistentHash
- RoundRobin
- Random

## ⚙️ Installation

```go
go get -u github.com/fufuok/load-balancer
```

## ⚡️ Quickstart

```go
package main

import (
	"fmt"

	balancer "github.com/fufuok/load-balancer"
)

func main() {
	choices := []*balancer.Choice{
		{Item: "🍒", Weight: 5},
		{Item: "🍋", Weight: 3},
		{Item: "🍉", Weight: 1},
		{Item: "🥑", Weight: 0},
	}
	balancer.Update(choices)

	// result of smooth selection is similar to: 🍒 🍒 🍒 🍋 🍒 🍋 🍒 🍋 🍉
	for i := 0; i < 9; i++ {
		fmt.Print(balancer.Select(), " ")
	}
	fmt.Println()
}
```

## 📚 Examples

please see: [examples](examples)

### Initialize the balancer

Sample data:

```go
var choices []*balancer.Choice

// for WeightedRoundRobin/SmoothWeightedRoundRobin/WeightedRand
// To be selected : Weighted
wNodes := map[string]int{
    "A": 5,
    "B": 3,
    "C": 1,
    "D": 0,
}
choices = balancer.NewChoicesMap(wNodes)

// or
choices = []*balancer.Choice{
    {Item: "A", Weight: 5},
    {Item: "B", Weight: 3},
    {Item: "C", Weight: 1},
    {Item: "D", Weight: 0},
}

// for RoundRobin/Random/ConsistentHash
nodes := []string{"A", "B", "C"}
choices = balancer.NewChoicesSlice(nodes)

// or
choices = []*balancer.Choice{
    {Item: "A"},
    {Item: "B"},
    {Item: "C"},
}
```

1. use default balancer (WRR)

   WeightedRoundRobin is the default balancer algorithm.

   ```go
   balancer.Update(choices)
   ```

2. use WeightedRoundRobin (WRR)

   ```go
   var lb balancer.Balancer
   lb = balancer.New(balancer.WeightedRoundRobin, choices)
   
   // or
   lb = balancer.New(balancer.WeightedRoundRobin, nil)
   lb.Update(choices)
   
   // or
   lb = balancer.NewWeightedRoundRobin(choices...)
   
   // or
   lb = balancer.NewWeightedRoundRobin()
   lb.Update(choices)
   ```

3. use SmoothWeightedRoundRobin (SWRR)

   ```go
   var lb balancer.Balancer
   lb = balancer.New(balancer.SmoothWeightedRoundRobin, choices)
   
   // or
   lb = balancer.New(balancer.SmoothWeightedRoundRobin, nil)
   lb.Update(choices)
   
   // or
   lb = balancer.NewSmoothWeightedRoundRobin(choices...)
   
   // or
   lb = balancer.NewSmoothWeightedRoundRobin()
   lb.Update(choices)
   ```

4. use WeightedRand (WR)

   ```go
   var lb balancer.Balancer
   lb = balancer.New(balancer.WeightedRand, choices)
   
   // or
   lb = balancer.New(balancer.WeightedRand, nil)
   lb.Update(choices)
   
   // or
   lb = balancer.NewWeightedRand(choices...)
   
   // or
   lb = balancer.NewWeightedRand()
   lb.Update(choices)
   ```

5. use ConsistentHash

   ```go
   var lb balancer.Balancer
   lb = balancer.New(balancer.ConsistentHash, choices)
   
   // or
   lb = balancer.New(balancer.ConsistentHash, nil)
   lb.Update(choices)
   
   // or
   lb = balancer.NewConsistentHash(choices...)
   
   // or
   lb = balancer.NewConsistentHash()
   lb.Update(choices)
   ```

6. use RoundRobin (RR)

   ```go
   var lb balancer.Balancer
   lb = balancer.New(balancer.RoundRobin, choices)
   
   // or
   lb = balancer.New(balancer.RoundRobin, nil)
   lb.Update(choices)
   
   // or
   lb = balancer.NewRoundRobin(choices...)
   
   // or
   lb = balancer.NewRoundRobin()
   lb.Update(choices)
   ```

7. use Random

   ```go
   var lb balancer.Balancer
   lb = balancer.New(balancer.Random, choices)
   
   // or
   lb = balancer.New(balancer.Random, nil)
   lb.Update(choices)
   
   // or
   lb = balancer.NewRandom(choices...)
   
   // or
   lb = balancer.NewRandom()
   lb.Update(choices)
   ```

### Gets next selected item

```go
node := lb.Select()
```

ip consistent hash:

```go
node := lb.Select("192.168.1.100")
node := lb.Select("192.168.1.100", "Test", "...")
```

### Interface

```go
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
```

## 🤖 Benchmarks

```shell
go test -run=^$ -benchmem -benchtime=1s -count=1 -bench=.
goos: linux
goarch: amd64
pkg: github.com/fufuok/load-balancer
cpu: Intel(R) Xeon(R) Gold 6151 CPU @ 3.00GHz
BenchmarkBalancer/WRR/10-4                              42144165                25.84 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR/10-4                             69472514                15.81 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR/10-4                               50838763                22.90 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash/10-4                             39170830                30.32 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin/10-4                       84562006                14.10 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random/10-4                          244497501                4.894 ns/op            0 B/op          0 allocs/op

BenchmarkBalancer/WRR#01/100-4                          41521372                30.54 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR#01/100-4                          7069575                170.0 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR#01/100-4                           31241767                39.10 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash#01/100-4                         31805023                37.79 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin#01/100-4                   84337191                14.24 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random#01/100-4                      245488258                4.919 ns/op            0 B/op          0 allocs/op

BenchmarkBalancer/WRR#02/1000-4                         34251942                34.42 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR#02/1000-4                          881721                 1368 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR#02/1000-4                          19458044                61.63 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash#02/1000-4                        23156967                51.16 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin#02/1000-4                  84485420                14.14 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random#02/1000-4                     232620525                5.089 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR/10-4                      17618566                68.15 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR/10-4                      5395550                227.7 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/WR/10-4                      212595266                6.172 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash/10-4                    150920656                8.051 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin/10-4               42911876                28.01 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random/10-4                  204149901                5.905 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR#01/100-4                  18904735                67.70 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR#01/100-4                  1308855                906.3 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/WR#01/100-4                  117712749                10.34 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash#01/100-4                122846144                12.29 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin#01/100-4           41949963                28.52 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random#01/100-4              239025558                6.297 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR#02/1000-4                 17478133                68.95 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR#02/1000-4                  442495                 2562 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/WR#02/1000-4                  74902462                16.06 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash#02/1000-4                80515016                13.98 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin#02/1000-4          40136712                28.31 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random#02/1000-4             197248365                6.224 ns/op            0 B/op          0 allocs/op
```

## ⚠️ License

Third-party library licenses:

- [doublejump]([doublejump/LICENSE at master · edwingeng/doublejump (github.com)](https://github.com/edwingeng/doublejump/blob/master/LICENSE))
- [go-jump]([go-jump/LICENSE at master · dgryski/go-jump (github.com)](https://github.com/dgryski/go-jump/blob/master/LICENSE))





*ff*