# ⚖️ load balancing algorithm library

High-performance general load balancing algorithm library, non-goroutine-safe.

Smooth weighted load balancing algorithm: [NGINX](https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35) and [LVS](http://kb.linuxvirtualserver.org/wiki/Weighted_Round-Robin_Scheduling), Doublejump provides a revamped Google's jump consistent hash.

---

If you want a **goroutine-safe** load balancer, please refer to another library, which supports more APIs: [fufuok/balancer](https://github.com/fufuok/balancer)

---

## 🎯 Features

- WeightedRoundRobin
- SmoothWeightedRoundRobin
- WeightedRand
- ConsistentHash
- RoundRobin
- Random

## ⚙️ Installation

```go
go get -u github.com/shibingli/load-balancer
```

## ⚡️ Quickstart

```go
package main

import (
	"fmt"

	balancer "github.com/shibingli/load-balancer"
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
BenchmarkBalancer/WRR-10-4                         46217316                25.75 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR-10-4                        77468209                15.27 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR-10-4                          52319631                22.49 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash-10-4                        39925784                30.21 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin-10-4                  66584161                18.14 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random-10-4                     242458752                4.904 ns/op            0 B/op          0 allocs/op

BenchmarkBalancer/WRR-100-4                        40734304                30.44 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR-100-4                        7942176                152.7 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR-100-4                         31176547                38.46 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash-100-4                       31901799                38.02 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin-100-4                 66206282                18.05 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random-100-4                    244788913                4.885 ns/op            0 B/op          0 allocs/op

BenchmarkBalancer/WRR-1000-4                       34538292                34.18 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR-1000-4                        843820                 1320 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR-1000-4                        20156709                59.92 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash-1000-4                      23270704                51.25 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin-1000-4                66487467                18.08 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random-1000-4                   221386104                5.098 ns/op            0 B/op          0 allocs/op

BenchmarkBalancer/WRR-10000-4                      34621371                34.29 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR-10000-4                        88591                14274 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR-10000-4                       14898531                81.10 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash-10000-4                     19238934                62.80 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin-10000-4               66228592                18.31 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random-10000-4                  191338644                6.872 ns/op            0 B/op          0 allocs/op

BenchmarkBalancer/WRR-100000-4                     30867440                34.63 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR-100000-4                        6897               173933 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR-100000-4                      10732960                111.6 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash-100000-4                    19251387                62.78 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin-100000-4              65002867                18.33 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random-100000-4                  88514517                13.07 ns/op            0 B/op          0 allocs/op

BenchmarkBalancer/WRR-1000000-4                    28559396                36.27 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/SWRR-1000000-4                        363              3264391 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/WR-1000000-4                      4875242                238.9 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Hash-1000000-4                   19132564                62.31 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/RoundRobin-1000000-4             64232452                18.44 ns/op            0 B/op          0 allocs/op
BenchmarkBalancer/Random-1000000-4                 25495011                45.96 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR-10-4                 18472707                67.69 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR-10-4                 4228370                283.4 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/WR-10-4                 208656978                6.016 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash-10-4               144690333                8.001 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin-10-4          19509270                61.33 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random-10-4             919527031                1.331 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR-100-4                19207489                66.60 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR-100-4                1311392                899.7 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/WR-100-4                100000000                10.26 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash-100-4              120191985                9.945 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin-100-4         18468038                62.90 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random-100-4            921913453                1.324 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR-1000-4               17837187                68.50 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR-1000-4                534634                 2184 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/WR-1000-4                78022376                15.89 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash-1000-4              87318837                13.83 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin-1000-4        18960835                61.85 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random-1000-4           214012765                5.450 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR-10000-4              17449467                65.24 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR-10000-4               164738                 7878 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/WR-10000-4               56703684                20.90 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash-10000-4             71464341                17.37 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin-10000-4       19946966                59.33 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random-10000-4          179242167                6.321 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR-100000-4             17843058                63.60 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR-100000-4               17827                66743 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/WR-100000-4              41696436                29.32 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash-100000-4            71783467                16.66 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin-100000-4      18777990                59.97 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random-100000-4         130029762                9.360 ns/op            0 B/op          0 allocs/op

BenchmarkBalancerParallel/WRR-1000000-4            12799609                78.48 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/SWRR-1000000-4               1447               853293 ns/op            1 B/op          0 allocs/op
BenchmarkBalancerParallel/WR-1000000-4             19860668                58.98 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Hash-1000000-4           72376710                16.90 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/RoundRobin-1000000-4     19127218                62.03 ns/op            0 B/op          0 allocs/op
BenchmarkBalancerParallel/Random-1000000-4         75088678                14.43 ns/op            0 B/op          0 allocs/op
```

## ⚠️ License

Third-party library licenses:

- [doublejump]([doublejump/LICENSE at master · edwingeng/doublejump (github.com)](https://github.com/edwingeng/doublejump/blob/master/LICENSE))
- [go-jump]([go-jump/LICENSE at master · dgryski/go-jump (github.com)](https://github.com/dgryski/go-jump/blob/master/LICENSE))

_ff_
