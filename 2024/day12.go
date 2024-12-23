package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day12.input", "Relative file path to use as input.")

type pair struct {
	r, c int
}

func (p pair) neighbours() [4]pair {
	return [4]pair{
		{p.r - 1, p.c}, // up
		{p.r, p.c + 1}, // right
		{p.r + 1, p.c}, // down
		{p.r, p.c - 1}, // left
	}
}

func (p pair) diags() [4]pair {
	return [4]pair{
		{p.r - 1, p.c + 1}, // up-right
		{p.r + 1, p.c + 1}, // down-right
		{p.r + 1, p.c - 1}, // down-left
		{p.r - 1, p.c - 1}, // up-left
	}
}

type garden map[pair]rune

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	board := make(garden)
	for r, s := range split[:len(split)-1] {
		for c, v := range s {
			board[pair{r, c}] = v
		}
	}

	plots := make(map[pair]int)
	areas := make(map[int]int)
	perimeters := make(map[int]int)
	corners := make(map[int]int)
	plot := 1
	for loc, kind := range board {
		q := []pair{loc}
		for len(q) > 0 {
			work := q[0]
			q = q[1:]
			if plots[work] != 0 {
				continue
			}
			plots[work] = plot
			areas[plot]++
			var different [4]bool
			for i, n := range work.neighbours() {
				if board[n] != kind {
					perimeters[plot]++
					different[i] = true
				} else {
					q = append(q, n)
				}
			}
			diags := work.diags()
			for i := range len(different) {
				// outer corners
				if different[i] && different[(i+1)%4] {
					corners[plot]++
				}
				// inner corners
				if !different[i] && !different[(i+1)%4] && board[diags[i]] != kind {
					corners[plot]++
				}
			}
		}
		plot++
	}

	var sumA, sumB int
	for k, area := range areas {
		sumA += area * perimeters[k]
		sumB += area * corners[k]
	}
	fmt.Println(sumA)
	fmt.Println(sumB)
}
