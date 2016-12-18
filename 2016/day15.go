package main

import (
	"fmt"
)

type Disc struct {
	Position, Positions uint
}

func (d Disc) Tick(delta uint) uint {
	return (d.Position + delta) % d.Positions
}

type DiscStack []Disc

func (ds DiscStack) Check(start uint) int {
	for i, d := range ds {
		if d.Tick(uint(i+1)+start) != 0 {
			return i
		}
	}
	return -1
}

func main() {
	discs := DiscStack{{0, 7}, {0, 13}, {2, 3}, {2, 5}, {0, 17}, {7, 19}, {0, 11}}
	inc := uint(1)
	current := 0
	for i := uint(0); ; i += inc {
		if fail := discs.Check(uint(i)); fail == -1 {
			fmt.Printf("Found match: %d\n", i)
			break
		} else if fail > current {
			for p := current; p < fail; p++ {
				inc *= discs[p].Positions
			}
			current = fail
		}
		fmt.Println(i)
	}
}
