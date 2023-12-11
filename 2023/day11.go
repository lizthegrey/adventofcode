package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day11.input", "Relative file path to use as input.")

type Coord struct {
	R, C int
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	stars := make(map[Coord]bool)
	bigRows := make(map[int]bool)
	bigCols := make(map[int]bool)

	var maxR, maxC int
	for r, line := range split[:len(split)-1] {
		maxR = r
		for c, v := range line {
			if c > maxC {
				maxC = c
			}
			if v == '#' {
				stars[Coord{r, c}] = true
			}
		}
	}

outerR:
	for r := 0; r <= maxR; r++ {
		for c := 0; c <= maxC; c++ {
			if stars[Coord{r, c}] {
				continue outerR
			}
		}
		bigRows[r] = true
	}
outerC:
	for c := 0; c <= maxC; c++ {
		for r := 0; r <= maxR; r++ {
			if stars[Coord{r, c}] {
				continue outerC
			}
		}
		bigCols[c] = true
	}

	var sumA, sumB int
	for x := range stars {
		for y := range stars {
			if x.R > y.R || x.R == y.R && x.C >= y.C {
				continue
			}
			for r := x.R; r != y.R; r++ {
				sumA++
				sumB++
				if bigRows[r] {
					sumA++
					sumB += 999999
				}
			}
			incr := 1
			if x.C > y.C {
				incr = -1
			}
			for c := x.C; c != y.C; c += incr {
				sumA++
				sumB++
				if bigCols[c] {
					sumA++
					sumB += 999999
				}
			}
		}
	}
	fmt.Println(sumA)
	fmt.Println(sumB)
}
