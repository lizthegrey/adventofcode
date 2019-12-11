package main

import (
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
	"math"
)

var inputFile = flag.String("inputFile", "inputs/day11.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use Part B logic.")

type Coord struct {
	X, Y int
}

func main() {
	flag.Parse()
	tape := intcode.ReadInput(*inputFile)
	if tape == nil {
		fmt.Println("Failed to parse input.")
		return
	}

	workingTape := tape.Copy()
	input := make(chan int, 1)
	output, done := workingTape.Process(input)

	panels := make(map[Coord]bool)
	loc := Coord{0, 0}
	// Up
	dir := 0

	if *partB {
		panels[loc] = true
	}

outer:
	for {
		if panels[loc] {
			input <- 1
		} else {
			input <- 0
		}
		select {
		case <-done:
			break outer
		case color := <-output:
			if color == 0 {
				panels[loc] = false
			} else {
				panels[loc] = true
			}
			turn := <-output
			if turn == 0 {
				dir--
			} else {
				dir++
			}
			if dir >= 4 {
				dir -= 4
			} else if dir < 0 {
				dir += 4
			}
			switch dir {
			// UP
			case 0:
				loc.Y++
			// RIGHT
			case 1:
				loc.X++
			// DOWN
			case 2:
				loc.Y--
			// LEFT
			case 3:
				loc.X--
			}
		}
	}
	fmt.Println(len(panels))

	minX := math.MaxInt32
	maxX := math.MinInt32
	minY := math.MaxInt32
	maxY := math.MinInt32
	for coord := range panels {
		if !panels[coord] {
			continue
		}
		if coord.X > maxX {
			maxX = coord.X
		}
		if coord.X < minX {
			minX = coord.X
		}
		if coord.Y > maxY {
			maxY = coord.Y
		}
		if coord.Y < minY {
			minY = coord.Y
		}
	}
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			if panels[Coord{x, y}] {
				fmt.Printf("x")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}
}
