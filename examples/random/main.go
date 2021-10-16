package main

import (
	"fmt"

	balancer "github.com/fufuok/load-balancer"
)

func main() {
	var choices []*balancer.Choice

	// for RoundRobin/Random/ConsistentHash
	nodes := []string{"A", "B", "C"}
	choices = balancer.NewChoicesSlice(nodes)

	// or
	// choices = []*balancer.Choice{
	// 	{Item: "A"},
	// 	{Item: "B"},
	// 	{Item: "C"},
	// }

	var lb balancer.Balancer
	lb = balancer.New(balancer.Random, choices)

	// or
	// lb = balancer.New(balancer.Random, nil)
	// lb.Update(choices)

	// or
	// lb = balancer.NewRandom(choices...)

	// or
	// lb = balancer.NewRandom()
	// lb.Update(choices)

	fmt.Println("balancer name:", lb.Name())

	// random
	for i := 0; i < 9; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()
}
