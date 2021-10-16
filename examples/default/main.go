package main

import (
	"fmt"

	balancer "github.com/fufuok/load-balancer"
)

func main() {
	// WeightedRoundRobin is the default balancer algorithm.
	fmt.Println("default balancer name:", balancer.Name())

	// reinitialize the balancer items.
	// D: will not select items with a weight of 0
	choices := []*balancer.Choice{
		{Item: "A", Weight: 5},
		{Item: "B", Weight: 3},
		{Item: "C", Weight: 1},
		{Item: "D", Weight: 0},
	}
	balancer.Update(choices)

	// result of smooth selection is similar to: A A A B A B A B C
	for i := 0; i < 9; i++ {
		fmt.Print(balancer.Select(), " ")
	}
	fmt.Println()
}
