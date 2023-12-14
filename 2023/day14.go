package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day14.input", "Relative file path to use as input.")

type Coord struct {
	R, C int
}

var North Coord = Coord{-1, 0}
var East Coord = Coord{0, 1}
var South Coord = Coord{1, 0}
var West Coord = Coord{0, -1}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var h uint64
	zobrist := make(map[Coord]uint64)

	stationary := make(map[Coord]bool)
	orig := make(map[Coord]bool)
	var maxR int
	maxC := len(split[0]) - 1
	for r, s := range split[:len(split)-1] {
		for c, v := range s {
			loc := Coord{r, c}
			val := rand.Uint64()
			zobrist[loc] = val
			switch v {
			case 'O':
				orig[loc] = true
				h ^= val
			case '#':
				stationary[loc] = true
			}
		}
		maxR = r
	}

	moving := make(map[Coord]bool)
	for loc := range orig {
		moving[loc] = true
	}
	step(North, maxR, maxC, stationary, moving, h, zobrist)
	fmt.Println(sum(moving, maxR))

	clear(moving)
	for loc := range orig {
		moving[loc] = true
	}

	dirs := []Coord{North, West, South, East}
	zCache := map[uint64]int{
		h: 0,
	}
	var loopLen int
	for i := 0; i < 1000000000; i++ {
		for j := 0; j < len(dirs); j++ {
			h = step(dirs[j%4], maxR, maxC, stationary, moving, h, zobrist)
		}
		// If we've seen this exact situation before, we can skip ahead a bunch.
		// This allows us to avoid literally running 1B times.
		if turn, ok := zCache[h]; ok && loopLen == 0 {
			loopLen = i - turn
			for i+loopLen < 1000000000 {
				i += loopLen
			}
		} else {
			zCache[h] = i
		}
	}
	fmt.Println(sum(moving, maxR))
}

func step(dir Coord, maxR, maxC int, stationary, moving map[Coord]bool, h uint64, zobrist map[Coord]uint64) uint64 {
	// iteratively move rocks until they can move no further.
	// There is no point optimising this to run 10 or 100x faster by doing
	// clever tricks involving one column at a time, because iterating even
	// those 1B times will be too expensive.
outer:
	for {
		for loc := range moving {
			moved := Coord{loc.R + dir.R, loc.C + dir.C}
			if moved.R < 0 || moved.C < 0 || moved.R > maxR || moved.C > maxC {
				continue
			}
			if stationary[moved] || moving[moved] {
				continue
			}
			h ^= zobrist[loc]
			h ^= zobrist[moved]
			moving[moved] = true
			delete(moving, loc)
			continue outer
		}
		// Nothing changed; we're done.
		return h
	}
}

func debug(maxR, maxC int, stationary, moving map[Coord]bool) {
	for r := 0; r <= maxR; r++ {
		for c := 0; c <= maxC; c++ {
			loc := Coord{r, c}
			if moving[loc] {
				fmt.Printf("O")
			} else if stationary[loc] {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func sum(moving map[Coord]bool, maxR int) int {
	var sum int
	for loc := range moving {
		sum += maxR - loc.R + 1
	}
	return sum
}
