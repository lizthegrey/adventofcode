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

	r := ring.New(1)
	Write(r, 0)

	for i := 1; i < *iterations+1; i++ {
		for s := 0; s < *steps; s++ {
			r = r.Next()
		}
		insert := ring.New(1)
		Write(insert, i)
		r = r.Link(insert).Prev()
	}
	fmt.Println(Read(r.Next()))
}

func Read(r *ring.Ring) int {
	return r.Value.(int)
}
func Write(r *ring.Ring, i int) {
	r.Value = i
}
