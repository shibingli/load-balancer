package balancer

import (
	"math/rand"
	"strconv"
	"testing"
)

const (
	numMin = 10
	numMax = 1_000_000
)

func BenchmarkBalancer(b *testing.B) {
	for n := numMin; n <= numMax; n *= 10 {
		b.Run("WRR-"+strconv.Itoa(n), func(b *testing.B) {
			choices := genChoices(n)
			lb := NewWeightedRoundRobin(choices...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				lb.Select()
			}
		})

		b.Run("SWRR-"+strconv.Itoa(n), func(b *testing.B) {
			choices := genChoices(n)
			lb := NewSmoothWeightedRoundRobin(choices...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				lb.Select()
			}
		})

		b.Run("WR-"+strconv.Itoa(n), func(b *testing.B) {
			choices := genChoices(n)
			lb := NewWeightedRand(choices...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				lb.Select()
			}
		})

		b.Run("Hash-"+strconv.Itoa(n), func(b *testing.B) {
			choices := genChoices(n)
			lb := NewConsistentHash(choices...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				lb.Select("192.168.1.1")
			}
		})

		b.Run("RoundRobin-"+strconv.Itoa(n), func(b *testing.B) {
			choices := genChoices(n)
			lb := NewRoundRobin(choices...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				lb.Select()
			}
		})

		b.Run("Random-"+strconv.Itoa(n), func(b *testing.B) {
			choices := genChoices(n)
			lb := NewRandom(choices...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				lb.Select()
			}
		})
	}
}

func BenchmarkBalancerParallel(b *testing.B) {
	for n := numMin; n <= numMax; n *= 10 {
		b.Run("WRR-"+strconv.Itoa(n), func(b *testing.B) {
			choices := genChoices(n)
			lb := NewWeightedRoundRobin(choices...)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					lb.Select()
				}
			})
		})

		b.Run("SWRR-"+strconv.Itoa(n), func(b *testing.B) {
			choices := genChoices(n)
			lb := NewSmoothWeightedRoundRobin(choices...)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					lb.Select()
				}
			})
		})

		b.Run("WR-"+strconv.Itoa(n), func(b *testing.B) {
			choices := genChoices(n)
			lb := NewWeightedRand(choices...)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					lb.Select()
				}
			})
		})

		b.Run("Hash-"+strconv.Itoa(n), func(b *testing.B) {
			choices := genChoices(n)
			lb := NewConsistentHash(choices...)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					lb.Select("192.168.1.1")
				}
			})
		})

		b.Run("RoundRobin-"+strconv.Itoa(n), func(b *testing.B) {
			choices := genChoices(n)
			lb := NewRoundRobin(choices...)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					lb.Select()
				}
			})
		})

		b.Run("Random-"+strconv.Itoa(n), func(b *testing.B) {
			choices := genChoices(n)
			lb := NewRandom(choices...)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					lb.Select()
				}
			})
		})
	}
}

func genChoices(n int) []*Choice {
	choices := make([]*Choice, n)
	for i := 0; i < n; i++ {
		choices[i] = &Choice{
			Item:   strconv.Itoa(i),
			Weight: rand.Intn(20),
		}
	}
	return choices
}

// go test -run=^$ -benchmem -benchtime=1s -count=1 -bench=.
// goos: linux
// goarch: amd64
// pkg: github.com/fufuok/load-balancer
// cpu: Intel(R) Xeon(R) Gold 6151 CPU @ 3.00GHz
// BenchmarkBalancer/WRR-10-4                         46217316                25.75 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR-10-4                        77468209                15.27 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR-10-4                          52319631                22.49 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash-10-4                        39925784                30.21 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin-10-4                  66584161                18.14 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random-10-4                     242458752                4.904 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR-100-4                        40734304                30.44 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR-100-4                        7942176                152.7 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR-100-4                         31176547                38.46 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash-100-4                       31901799                38.02 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin-100-4                 66206282                18.05 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random-100-4                    244788913                4.885 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR-1000-4                       34538292                34.18 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR-1000-4                        843820                 1320 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR-1000-4                        20156709                59.92 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash-1000-4                      23270704                51.25 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin-1000-4                66487467                18.08 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random-1000-4                   221386104                5.098 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR-10000-4                      34621371                34.29 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR-10000-4                        88591                14274 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR-10000-4                       14898531                81.10 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash-10000-4                     19238934                62.80 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin-10000-4               66228592                18.31 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random-10000-4                  191338644                6.872 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR-100000-4                     30867440                34.63 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR-100000-4                        6897               173933 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR-100000-4                      10732960                111.6 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash-100000-4                    19251387                62.78 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin-100000-4              65002867                18.33 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random-100000-4                  88514517                13.07 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR-1000000-4                    28559396                36.27 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR-1000000-4                        363              3264391 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR-1000000-4                      4875242                238.9 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash-1000000-4                   19132564                62.31 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin-1000000-4             64232452                18.44 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random-1000000-4                 25495011                45.96 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR-10-4                 18472707                67.69 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR-10-4                 4228370                283.4 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR-10-4                 208656978                6.016 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash-10-4               144690333                8.001 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin-10-4          19509270                61.33 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random-10-4             919527031                1.331 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR-100-4                19207489                66.60 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR-100-4                1311392                899.7 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR-100-4                100000000                10.26 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash-100-4              120191985                9.945 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin-100-4         18468038                62.90 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random-100-4            921913453                1.324 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR-1000-4               17837187                68.50 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR-1000-4                534634                 2184 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR-1000-4                78022376                15.89 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash-1000-4              87318837                13.83 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin-1000-4        18960835                61.85 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random-1000-4           214012765                5.450 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR-10000-4              17449467                65.24 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR-10000-4               164738                 7878 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR-10000-4               56703684                20.90 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash-10000-4             71464341                17.37 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin-10000-4       19946966                59.33 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random-10000-4          179242167                6.321 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR-100000-4             17843058                63.60 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR-100000-4               17827                66743 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR-100000-4              41696436                29.32 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash-100000-4            71783467                16.66 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin-100000-4      18777990                59.97 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random-100000-4         130029762                9.360 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR-1000000-4            12799609                78.48 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR-1000000-4               1447               853293 ns/op            1 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR-1000000-4             19860668                58.98 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash-1000000-4           72376710                16.90 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin-1000000-4     19127218                62.03 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random-1000000-4         75088678                14.43 ns/op            0 B/op          0 allocs/op
