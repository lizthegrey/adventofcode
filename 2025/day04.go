package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day04.input", "Relative file path to use as input.")

type Coord struct {
	R, C int
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)[:len(bytes)-1]
	roll := make(map[Coord]bool)

	var countA, countB int
	for r, s := range strings.Split(contents, "\n") {
		for c, v := range s {
			switch v {
			case '.':
				// nothing
			case '@':
				roll[Coord{r, c}] = true
			}
		}
	}
	for pass := 0; ; pass++ {
		mark := make(map[Coord]bool)
		for loc := range roll {
			var neigh int
			for rDelta := -1; rDelta <= 1; rDelta++ {
				for cDelta := -1; cDelta <= 1; cDelta++ {
					if roll[Coord{loc.R + rDelta, loc.C + cDelta}] {
						neigh++
					}
				}
			}
			// 4+1, including ourselves, which we already know to be occupied.
			if neigh < 5 {
				mark[loc] = true
			}
		}
		for loc := range mark {
			delete(roll, loc)
		}
		if pass == 0 {
			countA += len(mark)
		}
		countB += len(mark)
		if len(mark) == 0 {
			break
		}
	}
	fmt.Println(countA)
	fmt.Println(countB)
}
