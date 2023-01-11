package main

import (
	"fmt"

	balancer "github.com/shibingli/load-balancer"
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
	lb = balancer.New(balancer.SmoothWeightedRoundRobin, choices)

	// or
	// lb = balancer.New(balancer.SmoothWeightedRoundRobin, nil)
	// lb.Update(choices)

	// or
	// lb = balancer.NewSmoothWeightedRoundRobin(choices...)

	// or
	// lb = balancer.NewSmoothWeightedRoundRobin()
	// lb.Update(choices)

	fmt.Println("balancer name:", lb.Name())

	// result of smooth selection is: A B A C A B A
	for i := 0; i < 7; i++ {
		fmt.Print(lb.Select(), " ")
	}
	fmt.Println()
}
