package main

import (
	"container/ring"
	"flag"
	"fmt"
)

var steps = flag.Int("steps", 3, "The number of steps to take between inserts.")
var iterations = flag.Int("iterations", 2017, "The number of insertions to perform.")

func main() {
	flag.Parse()

	start := ring.New(1)
	Write(start, 0)

	r := start
	for i := 1; i < *iterations+1; i++ {
		for s := 0; s < *steps%r.Len(); s++ {
			r = r.Next()
		}
		insert := ring.New(1)
		Write(insert, i)
		r = r.Link(insert).Prev()
	}
	fmt.Printf("Value after 2017 is %d.\n", Read(r.Next()))

	fmt.Printf("Value after 0 is %d.\n", Read(start.Next()))
}

func Read(r *ring.Ring) int {
	return r.Value.(int)
}
func Write(r *ring.Ring, i int) {
	r.Value = i
}
