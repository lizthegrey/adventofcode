package main

import (
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
	"math"
	"math/rand"
)

var inputFile = flag.String("inputFile", "inputs/day15.input", "Relative file path to use as input.")

type Coord struct {
	X, Y int
}

func (c Coord) Move(dir int) Coord {
	ret := c
	switch dir {
	case 1:
		ret.Y++
	case 2:
		ret.Y--
	case 3:
		ret.X--
	case 4:
		ret.X++
	}
	return ret
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

	loc := Coord{0, 0}
	passable := map[Coord]bool{
		loc: true,
	}
	supply := Coord{0, 0}

	// Perform a random walk.
outer:
	for {
		direction := 1 + rand.Intn(4)
		// 1=N,2=S,3=W,4=E
		proposedLoc := loc.Move(direction)
		input <- direction
		select {
		case <-done:
			break outer
		case status := <-output:
			switch status {
			case 0:
				// Hit a wall. Mark tile impassable.
				passable[proposedLoc] = false
			case 1:
				// Moved. Mark tile passable.
				passable[proposedLoc] = true
				loc = proposedLoc
			case 2:
				// Found the oxygen supply.
				passable[proposedLoc] = true
				supply = proposedLoc
				break outer
			}
		}
	}

	// Print the maze out to check it.
	minX := math.MaxInt32
	maxX := math.MinInt32
	minY := math.MaxInt32
	maxY := math.MinInt32
	for coord := range passable {
		if !passable[coord] {
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
			c := Coord{x, y}
			if v, ok := passable[c]; ok {
				if v {
					if supply == c {
						fmt.Printf("S")
					} else {
						fmt.Printf(" ")
					}
				} else {
					fmt.Printf("x")
				}
			} else {
				fmt.Printf("?")
			}
		}
		fmt.Println()
	}

	// Perform a breadth-first search.
	start := Coord{0, 0}
	shortest := map[Coord]int{
		start: 0,
	}
	worklist := []Coord{start}
	for {
		w := worklist[0]
		for dir := 1; dir <= 4; dir++ {
			moved := w.Move(dir)
			if !passable[moved] {
				continue
			}
			if moved == supply {
				fmt.Println(shortest[w] + 1)
				return
			}
			if _, ok := shortest[moved]; ok {
				continue
			} else {
				shortest[moved] = shortest[w] + 1
				worklist = append(worklist, moved)
			}
		}
		worklist = worklist[1:]
	}
}
