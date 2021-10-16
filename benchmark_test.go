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
		b.Run("WRR", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := genChoices(n)
				lb := NewWeightedRoundRobin(choices...)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					lb.Select()
				}
			})
		})

		b.Run("SWRR", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := genChoices(n)
				lb := NewSmoothWeightedRoundRobin(choices...)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					lb.Select()
				}
			})
		})

		b.Run("WR", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := genChoices(n)
				lb := NewWeightedRand(choices...)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					lb.Select()
				}
			})
		})

		b.Run("Hash", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := genChoices(n)
				lb := NewConsistentHash(choices...)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					lb.Select("192.168.1.1")
				}
			})
		})

		b.Run("RoundRobin", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := genChoices(n)
				lb := NewRoundRobin(choices...)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					lb.Select()
				}
			})
		})

		b.Run("Random", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := genChoices(n)
				lb := NewRandom(choices...)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					lb.Select()
				}
			})
		})
	}
}

func BenchmarkBalancerParallel(b *testing.B) {
	for n := numMin; n <= numMax; n *= 10 {
		b.Run("WRR", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := genChoices(n)
				lb := NewWeightedRoundRobin(choices...)
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						lb.Select()
					}
				})
			})
		})

		b.Run("SWRR", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := genChoices(n)
				lb := NewSmoothWeightedRoundRobin(choices...)
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						lb.Select()
					}
				})
			})
		})

		b.Run("WR", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := genChoices(n)
				lb := NewWeightedRand(choices...)
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						lb.Select()
					}
				})
			})
		})

		b.Run("Hash", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := genChoices(n)
				lb := NewConsistentHash(choices...)
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						lb.Select("192.168.1.1")
					}
				})
			})
		})

		b.Run("RoundRobin", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := genChoices(n)
				lb := NewRoundRobin(choices...)
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						lb.Select()
					}
				})
			})
		})

		b.Run("Random", func(b *testing.B) {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := genChoices(n)
				lb := NewRandom(choices...)
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						lb.Select()
					}
				})
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
// BenchmarkBalancer/WRR/10-4                              42144165                25.84 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR/10-4                             69472514                15.81 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR/10-4                               50838763                22.90 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash/10-4                             39170830                30.32 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin/10-4                       84562006                14.10 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random/10-4                          244497501                4.894 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#01/100-4                          41521372                30.54 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#01/100-4                          7069575                170.0 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR#01/100-4                           31241767                39.10 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#01/100-4                         31805023                37.79 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#01/100-4                   84337191                14.24 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#01/100-4                      245488258                4.919 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#02/1000-4                         34251942                34.42 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#02/1000-4                          881721                 1368 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR#02/1000-4                          19458044                61.63 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#02/1000-4                        23156967                51.16 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#02/1000-4                  84485420                14.14 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#02/1000-4                     232620525                5.089 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#03/10000-4                        34250386                34.62 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#03/10000-4                          89247                13367 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR#03/10000-4                         14632105                82.15 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#03/10000-4                       19206060                62.69 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#03/10000-4                 84202992                14.34 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#03/10000-4                    190687468                6.350 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#04/100000-4                       30741852                35.02 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#04/100000-4                          6874               168835 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR#04/100000-4                        10501063                114.3 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#04/100000-4                      19062961                62.45 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#04/100000-4                84166652                14.15 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#04/100000-4                    90442380                12.92 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancer/WRR#05/1000000-4                      28341630                35.50 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/SWRR#05/1000000-4                          358              3273492 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/WR#05/1000000-4                        5127837                237.4 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Hash#05/1000000-4                     19134792                62.62 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/RoundRobin#05/1000000-4               83314928                14.55 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancer/Random#05/1000000-4                   24208584                42.61 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR/10-4                      17618566                68.15 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR/10-4                      5395550                227.7 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR/10-4                      212595266                6.172 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash/10-4                    150920656                8.051 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin/10-4               42911876                28.01 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random/10-4                  204149901                5.905 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#01/100-4                  18904735                67.70 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#01/100-4                  1308855                906.3 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR#01/100-4                  117712749                10.34 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#01/100-4                122846144                12.29 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#01/100-4           41949963                28.52 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#01/100-4              239025558                6.297 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#02/1000-4                 17478133                68.95 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#02/1000-4                  442495                 2562 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR#02/1000-4                  74902462                16.06 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#02/1000-4                80515016                13.98 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#02/1000-4          40136712                28.31 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#02/1000-4             197248365                6.224 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#03/10000-4                18417727                65.07 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#03/10000-4                 169377                 7086 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR#03/10000-4                 56400831                21.39 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#03/10000-4               68102100                18.25 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#03/10000-4         42354654                28.36 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#03/10000-4            194200598                6.835 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#04/100000-4               18112053                66.78 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#04/100000-4                 17480                68866 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR#04/100000-4                41026820                29.49 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#04/100000-4              72594279                17.00 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#04/100000-4        42061038                28.44 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#04/100000-4           121176079                9.707 ns/op            0 B/op          0 allocs/op
//
// BenchmarkBalancerParallel/WRR#05/1000000-4              18719948                70.35 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/SWRR#05/1000000-4                 1460               853531 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/WR#05/1000000-4               19019304                63.41 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Hash#05/1000000-4             70520688                16.77 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/RoundRobin#05/1000000-4       42368685                28.18 ns/op            0 B/op          0 allocs/op
// BenchmarkBalancerParallel/Random#05/1000000-4           76504468                13.82 ns/op            0 B/op          0 allocs/op
