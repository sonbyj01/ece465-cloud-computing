package main

import (
	"fmt"
	"time"
)

func f(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}

// speculatively color vertices
func coloring(Graph *g) {
	var nextMinValue = 0

	for _, node := range g.nodes {
		fmt.Println("DEBUGGING", index, "=>", element)

		for _, neighborNode := range element.nodes {
			if node.value == neighborNode.value {
				node.value++
			}
		}
	}
}

func main() {
	f("direct")

	go f("goroutine")

	go func(msg string) {
		fmt.Println(msg)
	}("going")

	time.Sleep(time.Second)
	fmt.Println("done")
}
