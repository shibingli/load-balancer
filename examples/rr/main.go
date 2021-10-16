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
	lb = balancer.New(balancer.RoundRobin, choices)

	// or
	// lb = balancer.New(balancer.RoundRobin, nil)
	// lb.Update(choices)

	// or
	// lb = balancer.NewRoundRobin(choices...)

	// or
	// lb = balancer.NewRoundRobin()
	// lb.Update(choices)

	fmt.Println("balancer name:", lb.Name())

	// A B C A B
	for i := 0; i < 5; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()
}
