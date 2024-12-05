package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day04.input", "Relative file path to use as input.")

const match = "XMAS"

type coord struct {
	r, c int
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	grid := make(map[coord]rune)
	var maxR, maxC int
	for r, row := range split[:len(split)-1] {
		if r > maxR {
			maxR = r
		}
		for c, v := range row {
			if c > maxC {
				maxC = c
			}
			grid[coord{r, c}] = v
		}
	}

	var sum int
	sum += checkAll(grid, maxR, maxC, coord{0, 1})
	sum += checkAll(grid, maxR, maxC, coord{0, -1})
	sum += checkAll(grid, maxR, maxC, coord{1, 0})
	sum += checkAll(grid, maxR, maxC, coord{-1, 0})
	sum += checkAll(grid, maxR, maxC, coord{1, 1})
	sum += checkAll(grid, maxR, maxC, coord{-1, -1})
	sum += checkAll(grid, maxR, maxC, coord{1, -1})
	sum += checkAll(grid, maxR, maxC, coord{-1, 1})
	fmt.Println(sum)

	var count int
	for r := 0; r <= maxR; r++ {
		for c := 0; c <= maxC; c++ {
			if grid[coord{r, c}] != 'A' {
				continue
			}
			// Check the diagonals
			if !(grid[coord{r + 1, c - 1}] == 'M' && grid[coord{r - 1, c + 1}] == 'S' || grid[coord{r + 1, c - 1}] == 'S' && grid[coord{r - 1, c + 1}] == 'M') {
				continue
			}
			if !(grid[coord{r + 1, c + 1}] == 'M' && grid[coord{r - 1, c - 1}] == 'S' || grid[coord{r + 1, c + 1}] == 'S' && grid[coord{r - 1, c - 1}] == 'M') {
				continue
			}
			count++
		}
	}
	fmt.Println(count)
}

func checkAll(grid map[coord]rune, maxR, maxC int, step coord) int {
	var count int
	for r := 0; r <= maxR; r++ {
		for c := 0; c <= maxC; c++ {
			if check(grid, coord{r, c}, step) {
				count++
			}
		}
	}
	return count
}

func check(grid map[coord]rune, start coord, step coord) bool {
	loc := start
	for _, v := range match {
		if v != grid[loc] {
			return false
		}
		loc.r += step.r
		loc.c += step.c
	}
	return true
}
