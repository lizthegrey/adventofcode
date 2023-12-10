package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day10.input", "Relative file path to use as input.")

type Coord struct {
	R, C int
}

type Direction int

const (
	Invalid = iota
	North
	East
	South
	West
)

func (c Coord) Travel(dir Direction) Coord {
	switch dir {
	case North:
		return Coord{c.R - 1, c.C}
	case East:
		return Coord{c.R, c.C + 1}
	case South:
		return Coord{c.R + 1, c.C}
	case West:
		return Coord{c.R, c.C - 1}
	}
	return Coord{-1, -1}
}

var Traversals = map[rune][5]Direction{
	//    Invalid  North    East     South    West
	'|': {Invalid, North, Invalid, South, Invalid},
	'-': {Invalid, Invalid, East, Invalid, West},
	'L': {Invalid, Invalid, Invalid, East, North},
	'J': {Invalid, Invalid, North, West, Invalid},
	'7': {Invalid, West, South, Invalid, Invalid},
	'F': {Invalid, East, Invalid, Invalid, South},
	'.': {Invalid, Invalid, Invalid, Invalid, Invalid},
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var start Coord

	grid := make(map[Coord]rune)
	for r, s := range split[:len(split)-1] {
		for c, v := range s {
			loc := Coord{r, c}
			grid[loc] = v
			if v == 'S' {
				start = loc
			}
		}
	}

	// Check each pipe around S, clockwise, until we find one that is connected
	cur := start
	var facing Direction
	boundary := make(map[Coord]bool)

	var worklist []Coord
	for _, dir := range []Direction{North, East, South, West} {
		next := start.Travel(dir)
		if fNext := Traversals[grid[next]][dir]; fNext != Invalid {
			facing = dir
			break
		}
	}

	for {
		boundary[cur] = true
		cur = cur.Travel(facing)
		if cur == start {
			break
		}

		newFacing := Traversals[grid[cur]][facing]
		// for facing (and newFacing if different): add exterior nodes.
		for _, f := range []Direction{facing, newFacing} {
			var dir Direction
			switch f {
			case North:
				dir = West
			case East:
				dir = North
			case South:
				dir = East
			case West:
				dir = South
			}
			worklist = append(worklist, cur.Travel(dir))
			if facing == newFacing {
				break
			}
		}
		facing = newFacing
	}
	fmt.Println(len(boundary) / 2)

	// Flood fill inward.
	seen := make(map[Coord]bool)
	inside := make(map[Coord]bool)
	for _, v := range worklist {
		seen[v] = true
	}
	for len(worklist) > 0 {
		item := worklist[0]
		worklist = worklist[1:]
		if boundary[item] {
			continue
		}
		inside[item] = true
		for _, dir := range []Direction{North, East, South, West} {
			next := item.Travel(dir)
			if !seen[next] && !boundary[next] {
				seen[next] = true
				worklist = append(worklist, next)
			}
		}
	}
	fmt.Println(len(inside))
}
