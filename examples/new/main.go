package main

import (
	"fmt"

	balancer "github.com/shibingli/load-balancer"
)

func main() {
	var lb balancer.Balancer

	// SmoothWeightedRoundRobin / WeightedRoundRobin / WeightedRand
	wNodes := map[string]int{
		"A": 5,
		"B": 1,
		"C": 1,
		"D": 0,
	}
	choices := balancer.NewChoicesMap(wNodes)
	lb = balancer.New(balancer.WeightedRoundRobin, choices)
	fmt.Println("balancer name:", lb.Name())

	lb = balancer.New(balancer.WeightedRand, choices)
	fmt.Println("balancer name:", lb.Name())

	lb = balancer.New(balancer.SmoothWeightedRoundRobin, choices)
	fmt.Println("balancer name:", lb.Name())

	// result of SmoothWeightedRoundRobin: A A B A C A A
	for i := 0; i < 7; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()

	// RoundRobin / Random / ConsistentHash
	nodes := []string{"A", "B", "C"}
	choices = balancer.NewChoicesSlice(nodes)
	lb = balancer.New(balancer.ConsistentHash, choices)
	fmt.Println("balancer name:", lb.Name())

	lb = balancer.New(balancer.Random, choices)
	fmt.Println("balancer name:", lb.Name())

	lb = balancer.New(balancer.RoundRobin, choices)
	fmt.Println("balancer name:", lb.Name())

	// result of RoundRobin: A B C A B C A
	for i := 0; i < 7; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()

	// same effect
	lb = balancer.New(balancer.RoundRobin, nil)
	lb.Update(choices)
	lb.Select()
}
