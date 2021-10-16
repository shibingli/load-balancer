package main

import (
	"fmt"

	balancer "github.com/fufuok/load-balancer"
)

func main() {
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
	// choices = []*balancer.Choice{
	// 	{Item: "A", Weight: 5},
	// 	{Item: "B", Weight: 3},
	// 	{Item: "C", Weight: 1},
	// 	{Item: "D", Weight: 0},
	// }

	var lb balancer.Balancer
	lb = balancer.New(balancer.WeightedRoundRobin, choices)

	// or
	// lb = balancer.New(balancer.WeightedRoundRobin, nil)
	// lb.Update(choices)

	// or
	// lb = balancer.NewWeightedRoundRobin(choices...)

	// or
	// lb = balancer.NewWeightedRoundRobin()
	// lb.Update(choices)

	fmt.Println("balancer name:", lb.Name())

	// result of smooth selection is similar to: A A A B A B A B C
	for i := 0; i < 9; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()
}
