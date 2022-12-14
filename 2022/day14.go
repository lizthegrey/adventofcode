package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day14.input", "Relative file path to use as input.")

type coord struct {
	x, y int
}

type grid map[coord]bool

func (c coord) occupied(terrain, sand grid) bool {
	return terrain[c] || sand[c]
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	terrain := make(grid)
	var maxY int
	for _, s := range split[:len(split)-1] {
		path := strings.Split(s, " -> ")
		var prev coord
		for i, elem := range path {
			parts := strings.Split(elem, ",")
			x, _ := strconv.Atoi(parts[0])
			y, _ := strconv.Atoi(parts[1])
			if y > maxY {
				maxY = y
			}
			loc := coord{x, y}
			terrain[loc] = true
			if i != 0 {
				// Fill in a straight line from prev to loc.
				if prev.x == x {
					incr := 1
					if prev.y > y {
						incr = -1
					}
					for iy := prev.y; iy != y; iy += incr {
						terrain[coord{x, iy}] = true
					}
				} else if prev.y == loc.y {
					incr := 1
					if prev.x > x {
						incr = -1
					}
					for ix := prev.x; ix != x; ix += incr {
						terrain[coord{ix, y}] = true
					}
				} else {
					fmt.Println("Asked to draw an invalid diagonal")
					return
				}
			}
			prev = loc
		}
	}

	// part A
	source := coord{500, 0}
	sand := make(grid)
outer:
	for {
		loc := source
	inner:
		for {
			nextY := loc.y + 1
			if nextY > maxY {
				break outer
			}
			down := coord{loc.x, nextY}
			left := coord{loc.x - 1, nextY}
			right := coord{loc.x + 1, nextY}
			for _, proposed := range [3]coord{down, left, right} {
				if !proposed.occupied(terrain, sand) {
					loc = proposed
					continue inner
				}
			}
			// Sand grain couldn't find a place to go and has come to rest.
			sand[loc] = true
			break
		}
	}
	fmt.Println(len(sand))

	// part B
	// Continue from part A, except now treat maxY+2 as impassable.
outerB:
	for {
		loc := source
	innerB:
		for {
			if source.occupied(terrain, sand) {
				break outerB
			}
			nextY := loc.y + 1
			if nextY == maxY+2 {
				// This can't fall any further.
				sand[loc] = true
				break
			}
			down := coord{loc.x, nextY}
			left := coord{loc.x - 1, nextY}
			right := coord{loc.x + 1, nextY}
			for _, proposed := range [3]coord{down, left, right} {
				if !proposed.occupied(terrain, sand) {
					loc = proposed
					continue innerB
				}
			}
			// Sand grain couldn't find a place to go and has come to rest.
			sand[loc] = true
			break
		}
	}
	fmt.Println(len(sand))
}
