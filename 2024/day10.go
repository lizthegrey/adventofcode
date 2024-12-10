package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day10.input", "Relative file path to use as input.")

type coord struct {
	r, c int
}

type grid map[coord]int

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	board := make(grid)
	for r, s := range split[:len(split)-1] {
		for c, v := range s {
			loc := coord{r, c}
			height := int(v - '0')
			board[loc] = height
		}
	}

	var sumA, sumB int
	for loc, height := range board {
		if height != 0 {
			continue
		}
		reachable := make(grid)
		rating := score(board, reachable, loc)
		sumA += len(reachable)
		sumB += rating
	}
	fmt.Println(sumA)
	fmt.Println(sumB)
}

func score(board, reachable grid, loc coord) int {
	height := board[loc]
	if height == 9 {
		reachable[loc] = 1
		return 1
	}
	toVisit := [4]coord{
		{loc.r + 1, loc.c},
		{loc.r - 1, loc.c},
		{loc.r, loc.c + 1},
		{loc.r, loc.c - 1},
	}
	var rating int
	for _, neigh := range toVisit {
		if height+1 == board[neigh] {
			rating += score(board, reachable, neigh)
		}
	}
	return rating
}
