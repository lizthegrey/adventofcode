package main

import (
	"flag"
	"fmt"
	"math"
)

var input = flag.Int("input", 277678, "The puzzle input value.")

type Coord struct {
	X, Y int  // Cartesian coordinates
}

type Heading int8

type Status struct {
	Loc Coord
	Dir Heading
}

func main() {
	flag.Parse()

	s := Status{Coord{}, 2}
	seen := make(map[Coord]bool)
	for i := 1; i < *input; i++ {
		seen[s.Loc] = true
		if r := s.RotateCCW().Travel(); seen[r.Loc] == false {
			s = r
		} else {
			s = s.Travel()
		}
	}
	fmt.Printf("Solution is %d (stored at (%d,%d))\n", int64(math.Abs(float64(s.Loc.X))+math.Abs(float64(s.Loc.Y))), s.Loc.X, s.Loc.Y)
}

func (s Status) Travel() Status {
	var c Coord
	switch s.Dir {
		case 0:  // North
			c = Coord{s.Loc.X, s.Loc.Y+1}
		case 1:  // East
			c = Coord{s.Loc.X+1, s.Loc.Y}
		case 2:  // South
			c = Coord{s.Loc.X, s.Loc.Y-1}
		case 3:  // West
			c = Coord{s.Loc.X-1, s.Loc.Y}
	}
	return Status{c, s.Dir}
}

func (s Status) RotateCCW() Status {
	return Status{s.Loc, (s.Dir+3)%4}
}
