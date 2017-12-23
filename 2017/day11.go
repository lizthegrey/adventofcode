package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day11.input", "Relative file path to use as input.")

//      \(0, 1)/
//       \    /        (2,  1)
//  (-1,0)+--+ (1,0)
//       /    \
//   ---+(0, 0)+-      (2,  0)
//       \    /
// (-1,-1)+--+ (1,-1)
//       /    \        (2, -1)
//      /(0,-1)\
type Coord struct {
	X, Y int
}

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}

	paths := strings.Split(string(bytes), "\n")
	for _, path := range paths {
		directions := strings.Split(path, ",")
		loc := Coord{}
		furthest := 0
		for _, dir := range directions {
			loc = loc.Move(dir)
			moves := loc.GoHome()
			if moves > furthest {
				furthest = moves
			}
		}
		fmt.Printf("At (%d,%d).\n", loc.X, loc.Y)
		moves := loc.GoHome()
		fmt.Printf("\nGot home in %d moves (furthest: %d).\n", moves, furthest)
	}
}

func (c Coord) GoHome() int {
	loc := c
	moves := 0
	for (loc != Coord{}) {
		dir := ""
		evenX := (loc.X%2 == 0)
		if loc.X > 0 {
			if loc.Y > 0 {
				dir = "sw"
			} else if loc.Y < 0 {
				dir = "nw"
			} else {
				if evenX {
					dir = "nw"
				} else {
					dir = "sw"
				}
			}
		} else if loc.X < 0 {
			if loc.Y > 0 {
				dir = "se"
			} else if loc.Y < 0 {
				dir = "ne"
			} else {
				if evenX {
					dir = "ne"
				} else {
					dir = "se"
				}
			}
		} else {
			// loc.X is 0
			if loc.Y < 0 {
				dir = "n"
			} else {
				dir = "s"
			}
		}
		loc = loc.Move(dir)
		moves++
	}
	return moves
}

func (c Coord) Move(dir string) Coord {
	ret := c
	evenX := (c.X%2 == 0)
	switch dir {
	case "n":
		ret.Y++
	case "ne":
		ret.X++
		if !evenX {
			ret.Y++
		}
	case "se":
		ret.X++
		if evenX {
			ret.Y--
		}
	case "s":
		ret.Y--
	case "sw":
		ret.X--
		if evenX {
			ret.Y--
		}
	case "nw":
		ret.X--
		if !evenX {
			ret.Y++
		}
	}
	return ret
}
