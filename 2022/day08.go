package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day08.input", "Relative file path to use as input.")

type tree struct {
	height  int
	visible bool
}

type forest [][]tree

func (f forest) check(x, y int, highest *int) bool {
	height := f[y][x].height
	if height > *highest {
		*highest = height
		return true
	}
	return false
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var grid forest

	for _, s := range split[:len(split)-1] {
		row := make([]tree, 0)
		for _, v := range s {
			d, _ := strconv.Atoi(string(v))
			row = append(row, tree{height: d})
		}
		grid = append(grid, row)
	}

	// part A
	for y := 0; y < len(grid); y++ {
		highest := -1
		for x := 0; x < len(grid); x++ {
			if grid.check(x, y, &highest) {
				grid[y][x].visible = true
			}
		}
		highest = -1
		for x := len(grid) - 1; x >= 0; x-- {
			if grid.check(x, y, &highest) {
				grid[y][x].visible = true
			}
		}
	}
	for x := 0; x < len(grid); x++ {
		highest := -1
		for y := 0; y < len(grid); y++ {
			if grid.check(x, y, &highest) {
				grid[y][x].visible = true
			}
		}
		highest = -1
		for y := len(grid) - 1; y >= 0; y-- {
			if grid.check(x, y, &highest) {
				grid[y][x].visible = true
			}
		}
	}
	var total int
	for _, row := range grid {
		for _, t := range row {
			if t.visible {
				total++
			}
		}
	}
	fmt.Println(total)

	// part B
	var best int
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid); c++ {
			self := grid[r][c].height

			highest := -1
			var left, right, up, down int
			for x := c + 1; x < len(grid); x++ {
				grid.check(x, r, &highest)
				if highest >= self || x == len(grid)-1 {
					right = x - c
					break
				}
			}
			highest = -1
			for x := c - 1; x >= 0; x-- {
				grid.check(x, r, &highest)
				if highest >= self || x == 0 {
					left = c - x
					break
				}
			}
			highest = -1
			for y := r + 1; y < len(grid); y++ {
				grid.check(c, y, &highest)
				if highest >= self || y == len(grid)-1 {
					down = y - r
					break
				}
			}
			highest = -1
			for y := r - 1; y >= 0; y-- {
				grid.check(c, y, &highest)
				if highest >= self || y == 0 {
					up = r - y
					break
				}
			}

			// Compute the score.
			score := up * down * left * right
			if score > best {
				best = score
			}
		}
	}
	fmt.Println(best)
}
