package main

import (
	"flag"
	"fmt"
	"math"
)

var input = flag.Int("input", 277678, "The puzzle input value.")
var partB = flag.Bool("partB", true, "Whether to use part B logic.")

type Coord struct {
	X, Y int // Cartesian coordinates
}

type Heading int8

type Status struct {
	Loc Coord
	Dir Heading
}

type Board map[Coord]uint64

func main() {
	flag.Parse()

	s := Status{Coord{}, 2}
	seen := make(Board)
	for i := 1; i < *input; i++ {
		if i == 1 {
			seen[s.Loc] = 1
		} else {
			seen[s.Loc] = seen.SumNeighbors(s.Loc)
			if *partB && seen[s.Loc] > uint64(*input) {
				break
			}
		}

		if r := s.RotateCCW().Travel(); seen[r.Loc] == 0 {
			s = r
		} else {
			s = s.Travel()
		}
	}
	seen[s.Loc] = seen.SumNeighbors(s.Loc)
	fmt.Printf("%d steps (%d stored at (%d,%d))\n", int64(math.Abs(float64(s.Loc.X))+math.Abs(float64(s.Loc.Y))), seen[s.Loc], s.Loc.X, s.Loc.Y)
}

func (b Board) SumNeighbors(c Coord) uint64 {
	r := b[Coord{c.X + 1, c.Y + 0}]
	r += b[Coord{c.X + 1, c.Y + 1}]
	r += b[Coord{c.X + 0, c.Y + 1}]
	r += b[Coord{c.X - 1, c.Y + 1}]
	r += b[Coord{c.X - 1, c.Y + 0}]
	r += b[Coord{c.X - 1, c.Y - 1}]
	r += b[Coord{c.X + 0, c.Y - 1}]
	r += b[Coord{c.X + 1, c.Y - 1}]

	return r
}

func (s Status) Travel() Status {
	var c Coord
	switch s.Dir {
	case 0: // North
		c = Coord{s.Loc.X, s.Loc.Y + 1}
	case 1: // East
		c = Coord{s.Loc.X + 1, s.Loc.Y}
	case 2: // South
		c = Coord{s.Loc.X, s.Loc.Y - 1}
	case 3: // West
		c = Coord{s.Loc.X - 1, s.Loc.Y}
	}
	return Status{c, s.Dir}
}

func (s Status) RotateCCW() Status {
	return Status{s.Loc, (s.Dir + 3) % 4}
}
