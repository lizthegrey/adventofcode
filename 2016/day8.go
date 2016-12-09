package main

import (
	"fmt"
	"strconv"
	"strings"
)

type coord struct {
	X, Y int
}

func main() {
	input := strings.Split(`rect 1x1
rotate row y=0 by 20
rect 1x1
rotate row y=0 by 2
rect 1x1
rotate row y=0 by 3
rect 2x1
rotate row y=0 by 2
rect 1x1
rotate row y=0 by 3
rect 2x1
rotate row y=0 by 2
rect 1x1
rotate row y=0 by 4
rect 2x1
rotate row y=0 by 2
rect 1x1
rotate row y=0 by 2
rect 1x1
rotate row y=0 by 2
rect 1x1
rotate row y=0 by 3
rect 2x1
rotate row y=0 by 2
rect 1x1
rotate row y=0 by 5
rect 1x1
rotate row y=0 by 2
rect 1x1
rotate row y=0 by 6
rect 5x1
rotate row y=0 by 2
rect 1x3
rotate row y=2 by 8
rotate row y=0 by 8
rotate column x=0 by 1
rect 7x1
rotate row y=2 by 24
rotate row y=0 by 20
rotate column x=5 by 1
rotate column x=4 by 2
rotate column x=2 by 2
rotate column x=0 by 1
rect 7x1
rotate column x=34 by 2
rotate column x=22 by 1
rotate column x=15 by 1
rotate row y=2 by 18
rotate row y=0 by 12
rotate column x=8 by 2
rotate column x=7 by 1
rotate column x=5 by 2
rotate column x=2 by 1
rotate column x=0 by 1
rect 9x1
rotate row y=3 by 28
rotate row y=1 by 28
rotate row y=0 by 20
rotate column x=18 by 1
rotate column x=15 by 1
rotate column x=14 by 1
rotate column x=13 by 1
rotate column x=12 by 2
rotate column x=10 by 3
rotate column x=8 by 1
rotate column x=7 by 2
rotate column x=6 by 1
rotate column x=5 by 1
rotate column x=3 by 1
rotate column x=2 by 2
rotate column x=0 by 1
rect 19x1
rotate column x=34 by 2
rotate column x=24 by 1
rotate column x=23 by 1
rotate column x=14 by 1
rotate column x=9 by 2
rotate column x=4 by 2
rotate row y=3 by 5
rotate row y=2 by 3
rotate row y=1 by 7
rotate row y=0 by 5
rotate column x=0 by 2
rect 3x2
rotate column x=16 by 2
rotate row y=3 by 27
rotate row y=2 by 5
rotate row y=0 by 20
rotate column x=8 by 2
rotate column x=7 by 1
rotate column x=5 by 1
rotate column x=3 by 3
rotate column x=2 by 1
rotate column x=1 by 2
rotate column x=0 by 1
rect 9x1
rotate row y=4 by 42
rotate row y=3 by 40
rotate row y=1 by 30
rotate row y=0 by 40
rotate column x=37 by 2
rotate column x=36 by 3
rotate column x=35 by 1
rotate column x=33 by 1
rotate column x=32 by 1
rotate column x=31 by 3
rotate column x=30 by 1
rotate column x=28 by 1
rotate column x=27 by 1
rotate column x=25 by 1
rotate column x=23 by 3
rotate column x=22 by 1
rotate column x=21 by 1
rotate column x=20 by 1
rotate column x=18 by 1
rotate column x=17 by 1
rotate column x=16 by 3
rotate column x=15 by 1
rotate column x=13 by 1
rotate column x=12 by 1
rotate column x=11 by 2
rotate column x=10 by 1
rotate column x=8 by 1
rotate column x=7 by 2
rotate column x=5 by 1
rotate column x=3 by 3
rotate column x=2 by 1
rotate column x=1 by 1
rotate column x=0 by 1
rect 39x1
rotate column x=44 by 2
rotate column x=42 by 2
rotate column x=35 by 5
rotate column x=34 by 2
rotate column x=32 by 2
rotate column x=29 by 2
rotate column x=25 by 5
rotate column x=24 by 2
rotate column x=19 by 2
rotate column x=15 by 4
rotate column x=14 by 2
rotate column x=12 by 3
rotate column x=9 by 2
rotate column x=5 by 5
rotate column x=4 by 2
rotate row y=5 by 5
rotate row y=4 by 38
rotate row y=3 by 10
rotate row y=2 by 46
rotate row y=1 by 10
rotate column x=48 by 4
rotate column x=47 by 3
rotate column x=46 by 3
rotate column x=45 by 1
rotate column x=43 by 1
rotate column x=37 by 5
rotate column x=36 by 5
rotate column x=35 by 4
rotate column x=33 by 1
rotate column x=32 by 5
rotate column x=31 by 5
rotate column x=28 by 5
rotate column x=27 by 5
rotate column x=26 by 3
rotate column x=25 by 4
rotate column x=23 by 1
rotate column x=17 by 5
rotate column x=16 by 5
rotate column x=13 by 1
rotate column x=12 by 5
rotate column x=11 by 5
rotate column x=3 by 1
rotate column x=0 by 1`, "\n")

	grid := make(map[coord]bool)
	gWidth := 50
	gHeight := 6
	for _, inst := range input {
		parsed := strings.Split(inst, " ")
		switch parsed[0] {
		case "rotate":
			amount, err := strconv.Atoi(parsed[4])
			if err != nil {
				fmt.Printf("Bad instruction: %s", parsed)
			}
			switch parsed[1] {
			case "column":
				x, err := strconv.Atoi(parsed[2][2:])
				if err != nil {
					fmt.Printf("Bad instruction: %s", parsed)
				}
				newCol := make([]bool, 0)
				for y := 0; y < gHeight; y++ {
					newCol = append(newCol, grid[coord{x, (y - amount + gHeight) % gHeight}])
				}
				for y := 0; y < gHeight; y++ {
					if newCol[y] {
						grid[coord{x, y}] = true
					} else {
						delete(grid, coord{x, y})
					}
				}

			case "row":
				y, err := strconv.Atoi(parsed[2][2:])
				if err != nil {
					fmt.Printf("Bad instruction: %s", parsed)
				}
				newRow := make([]bool, 0)
				for x := 0; x < gWidth; x++ {
					newRow = append(newRow, grid[coord{(x - amount + gWidth) % gWidth, y}])
				}
				for x := 0; x < gWidth; x++ {
					if newRow[x] {
						grid[coord{x, y}] = true
					} else {
						delete(grid, coord{x, y})
					}
				}
			}
		case "rect":
			dims := strings.Split(parsed[1], "x")
			w, err := strconv.Atoi(dims[0])
			if err != nil {
				fmt.Printf("Bad instruction: %s", parsed)
			}
			h, err := strconv.Atoi(dims[1])
			if err != nil {
				fmt.Printf("Bad instruction: %s", parsed)
			}
			for x := 0; x < w; x++ {
				for y := 0; y < h; y++ {
					grid[coord{x, y}] = true
				}
			}
		default:
			fmt.Printf("Bad instruction: %s", parsed)
			return
		}
	}
	fmt.Println(len(grid))
	for y := 0; y < gHeight; y++ {
		for x := 0; x < gWidth; x++ {
			if grid[coord{x, y}] {
				fmt.Printf("*")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}
}
