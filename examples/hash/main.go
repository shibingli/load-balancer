package main

import (
	"fmt"

	balancer "github.com/shibingli/load-balancer"
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
	lb = balancer.New(balancer.ConsistentHash, choices)

	// or
	// lb = balancer.New(balancer.ConsistentHash, nil)
	// lb.Update(choices)

	// or
	// lb = balancer.NewConsistentHash(choices...)

	// or
	// lb = balancer.NewConsistentHash()
	// lb.Update(choices)

	fmt.Println("balancer name:", lb.Name())

	// C C C C C
	for i := 0; i < 5; i++ {
		fmt.Print(lb.Select("192.168.1.1"), " ")
	}
	fmt.Println()
}
