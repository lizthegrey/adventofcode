package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day20.input", "Relative file path to use as input.")

type coord struct {
	r, c int
}

func (c coord) manhattan(o coord) int {
	return max(c.r, o.r) - min(c.r, o.r) + max(c.c, o.c) - min(c.c, o.c)
}

func (c coord) neigh() []coord {
	return []coord{
		{c.r + 1, c.c},
		{c.r - 1, c.c},
		{c.r, c.c + 1},
		{c.r, c.c - 1},
	}
}

func (c coord) landings(skip int) []coord {
	var ret []coord
	for rDelta := -skip; rDelta <= skip; rDelta++ {
		rAbs := max(rDelta, -rDelta)
		for cDelta := -skip; cDelta <= skip; cDelta++ {
			cAbs := max(cDelta, -cDelta)
			if rAbs+cAbs <= skip {
				ret = append(ret, coord{c.r + rDelta, c.c + cDelta})
			}
		}
	}
	return ret
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var start, end coord
	passable := make(map[coord]bool)
	for r, s := range split[:len(split)-1] {
		for c, v := range s {
			loc := coord{r, c}
			switch v {
			case '.':
				passable[loc] = true
			case 'S':
				passable[loc] = true
				start = loc
			case 'E':
				passable[loc] = true
				end = loc
			}
		}
	}
	fromStart := flood(passable, start)
	fromEnd := flood(passable, end)
	fmt.Println(shortcuts(passable, fromStart, fromEnd, end, 2))
	fmt.Println(shortcuts(passable, fromStart, fromEnd, end, 20))
}

func flood(passable map[coord]bool, origin coord) map[coord]int {
	ret := make(map[coord]int)
	var distance int
	work := map[coord]bool{
		origin: true,
	}
	for len(work) != 0 {
		next := make(map[coord]bool)
		for loc := range work {
			ret[loc] = distance
			for _, n := range loc.neigh() {
				if _, ok := ret[n]; !ok && passable[n] {
					next[n] = true
				}
			}
		}
		work = next
		distance++
	}
	return ret
}

func shortcuts(passable map[coord]bool, fromStart, fromEnd map[coord]int, end coord, skip int) uint64 {
	withoutCheating := fromStart[end]
	// Compute all valid land squares within manhattan distance skip of each valid jump point
	// Determine if sum of fill from start to jump, manhattan from jump to land, fill from land to end, is 100 or more shorter than shortest path
	var count uint64
	for jump := range passable {
		a := fromStart[jump]
		for _, land := range jump.landings(skip) {
			if !passable[land] {
				continue
			}
			dist := a + jump.manhattan(land) + fromEnd[land]
			if dist+100 <= withoutCheating {
				count++
			}
		}
	}
	return count
}
