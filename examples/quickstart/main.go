package main

import (
	"fmt"

	balancer "github.com/fufuok/load-balancer"
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
