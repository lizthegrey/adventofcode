package main

import (
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
	"math"
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
	output, _ := workingTape.Process(input)

	loc := Coord{0, 0}
	passable := map[Coord]bool{
		loc: true,
	}
	supply := Coord{0, 0}

	// Keep track of what's unexplored, and backtrack to try to find unexplored nodes.
	// Use a DFS in order to avoid repeated backtracking.
	unexplored := make([]Coord, 0)
	for d := 1; d <= 4; d++ {
		unexplored = append(unexplored, loc.Move(d))
	}

	for {
		toExplore := unexplored[0]
		if _, ok := passable[toExplore]; ok {
			// Short circuit.
			unexplored = unexplored[1:]
			continue
		}
		result := path(input, output, &loc, toExplore, passable)
		switch result {
		case 0:
			// Don't do anything; we have no unexplored nodes to add.
		case 2:
			supply = loc
			fallthrough
		case 1:
			// Add new nodes to explore.
			for d := 1; d <= 4; d++ {
				proposed := loc.Move(d)
				if _, ok := passable[proposed]; !ok {
					unexplored = append(unexplored, proposed)
				}
			}
		}

		if len(unexplored) == 1 {
			// Finished exploring the maze.
			break
		}
		unexplored = unexplored[1:]
	}

	fmt.Println("Finished exploring maze. Result:")
	start := Coord{0, 0}
	printMaze(passable, start, supply)
	fmt.Println(len(bfs(start, supply, passable)))

	fmt.Println(fill(supply, passable))
}

func printMaze(passable map[Coord]bool, loc, supply Coord) {
	// Print the maze out to check it.
	minX := math.MaxInt32
	maxX := math.MinInt32
	minY := math.MaxInt32
	maxY := math.MinInt32
	for coord := range passable {
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
						if loc == supply {
							fmt.Printf("S")
						} else {
							fmt.Printf("s")
						}
					} else {
						if loc == c {
							fmt.Printf("*")
						} else {
							fmt.Printf(" ")
						}
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
}

func bfs(src, dst Coord, passable map[Coord]bool) []int {
	// Perform a breadth-first search.
	shortest := map[Coord][]int{
		src: []int{},
	}
	worklist := []Coord{src}
	for {
		w := worklist[0]
		for dir := 1; dir <= 4; dir++ {
			moved := w.Move(dir)
			if !passable[moved] && moved != dst {
				// Allow ourselves to pass on unknown ground only for the last move.
				continue
			}
			if _, ok := shortest[moved]; ok {
				continue
			} else {
				directions := make([]int, len(shortest[w])+1)
				copy(directions, shortest[w])
				directions[len(shortest[w])] = dir
				shortest[moved] = directions
				if moved == dst {
					return shortest[moved]
				}
				worklist = append(worklist, moved)
			}
		}
		worklist = worklist[1:]
	}
}

func fill(src Coord, passable map[Coord]bool) int {
	// Perform a breadth-first search.
	shortest := map[Coord]int{
		src: 0,
	}

	worklist := []Coord{src}
	for {
		w := worklist[0]
		for dir := 1; dir <= 4; dir++ {
			moved := w.Move(dir)
			if !passable[moved] {
				continue
			}
			if _, ok := shortest[moved]; ok {
				continue
			} else {
				shortest[moved] = shortest[w] + 1
				worklist = append(worklist, moved)
			}
		}
		if len(worklist) == 1 {
			return shortest[w]
		}
		worklist = worklist[1:]
	}
}

func path(input, output chan int, loc *Coord, dst Coord, passable map[Coord]bool) int {
	toMove := bfs(*loc, dst, passable)
	lastStatus := -1

	for i, d := range toMove {
		// 1=N,2=S,3=W,4=E
		proposedLoc := loc.Move(d)
		input <- d
		status := <-output
		lastStatus = status
		switch status {
		case 0:
			// Hit a wall. Mark tile impassable.
			if i != len(toMove)-1 {
				// This should never happen.
				fmt.Println("Failed to traverse known path.")
			}
			passable[proposedLoc] = false
		case 1:
			fallthrough
		case 2:
			// Moved. Mark tile passable.
			// Includes the found the oxygen supply case.
			passable[proposedLoc] = true
			*loc = proposedLoc
		}
	}
	return lastStatus
}
