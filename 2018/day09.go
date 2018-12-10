package main

import (
	"container/ring"
	"flag"
	"fmt"
)

var numPlayers = flag.Int("numPlayers", 9, "The number of players.")
var maxMarble = flag.Int("maxMarble", 25, "The maximum marble score.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

func main() {
	flag.Parse()

	scores := make([]int, *numPlayers)
	pidx := -1
	current := ring.New(1)
	current.Value = 0

	for i := 1; i <= *maxMarble; i++ {
		pidx = (pidx + 1) % *numPlayers
		if i%23 != 0 {
			current = current.Next()
			added := ring.New(1)
			added.Value = i
			current.Link(added)
			current = current.Next()
		} else {
			scores[pidx] += i
			for v := 0; v < 7; v++ {
				current = current.Prev()
			}
			current = current.Prev()
			removed := current.Unlink(1)
			current = current.Next()
			scores[pidx] += removed.Value.(int)
		}
	}

	result := 0
	for _, v := range scores {
		if v > result {
			result = v
		}
	}
	fmt.Printf("Result is %d\n", result)
}
