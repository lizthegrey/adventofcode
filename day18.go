package main

import (
	"bufio"
	"fmt"
	"os"
)

type Coord struct {
	x, y int
}

type Grid map[Coord]int

func main() {
	reader := bufio.NewReader(os.Stdin)
	on := make(Grid)
	y := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		for x := range line[:len(line) - 1] {
			if line[x] == '#' {
				on[Coord{x, y}] = 1
			}
		}
		y++
	}

	on[Coord{0, 0}] = 1
	on[Coord{0, 99}] = 1
	on[Coord{99, 0}] = 1
	on[Coord{99, 99}] = 1

	for i := 0; i < 100; i++ {
		fmt.Println(len(on))
		on = step(on)
	}
	fmt.Println(len(on))
}

func step(on Grid) Grid {
	next := make(Grid)

	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			neighborsOn := on[Coord{x-1, y-1}] + on[Coord{x-1, y}] + on[Coord{x-1, y+1}] + on[Coord{x, y-1}] + on[Coord{x, y+1}] + on[Coord{x+1, y-1}] + on[Coord{x+1, y}] + on[Coord{x+1, y+1}]

			if on[Coord{x, y}] == 1 && (neighborsOn == 2 || neighborsOn == 3) {
				next[Coord{x, y}] = 1
			} else  if on[Coord{x, y}] == 0 && neighborsOn == 3 {
				next[Coord{x, y}] = 1
			}
		}
	}

	next[Coord{0, 0}] = 1
	next[Coord{0, 99}] = 1
	next[Coord{99, 0}] = 1
	next[Coord{99, 99}] = 1

	return next
}
