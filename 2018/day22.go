package main

import (
	"flag"
	"fmt"
)

var depth = flag.Int("depth", 8112, "The depth of the cave system.")
var targetX = flag.Int("targetX", 13, "The target's X coordinate.")
var targetY = flag.Int("targetY", 743, "The target's Y coordinate.")
var margin = flag.Int("margin", 150, "The margin by which we'll check 'longer' paths.")

type Position struct {
	X, Y  int
	Climb bool
	Torch bool
}

func (p Position) Passable(erosion [][]int) bool {
	if p.X < 0 || p.X >= len(erosion) {
		return false
	}
	if p.Y < 0 || p.Y >= len(erosion[0]) {
		return false
	}

	terrain := erosion[p.X][p.Y] % 3
	switch terrain {
	case 0:
		// Rocky
		return p.Climb || p.Torch
	case 1:
		// Wet
		return !p.Torch
	case 2:
		// Narrow
		return !p.Climb
	}
	return false
}

func (p Position) Moves(erosion [][]int) map[Position]int {
	ret := make(map[Position]int)
	terrain := erosion[p.X][p.Y] % 3
	swapGear := p
	switch terrain {
	case 0:
		// Rocky
		if !p.Climb && !p.Torch {
			// Illegal position.
			return ret
		}
		swapGear.Climb = !swapGear.Climb
		swapGear.Torch = !swapGear.Torch
	case 1:
		// Wet
		if swapGear.Torch {
			// Illegal position.
			return ret
		}
		swapGear.Climb = !swapGear.Climb
	case 2:
		// Narrow
		if swapGear.Climb {
			// Illegal position.
			return ret
		}
		swapGear.Torch = !swapGear.Torch
	}
	ret[swapGear] = 7

	up := p
	up.Y--
	if up.Passable(erosion) {
		ret[up] = 1
	}

	down := p
	down.Y++
	if down.Passable(erosion) {
		ret[down] = 1
	}

	left := p
	left.X--
	if left.Passable(erosion) {
		ret[left] = 1
	}

	right := p
	right.X++
	if right.Passable(erosion) {
		ret[right] = 1
	}
	return ret
}

func main() {
	flag.Parse()

	risk := 0
	// [x][y] in ordering.
	geoIdx := make([][]int, *targetX+1+*margin)
	erosion := make([][]int, *targetX+1+*margin)
	for x := 0; x <= *targetX+*margin; x++ {
		geoIdx[x] = make([]int, *targetY+1+*margin)
		erosion[x] = make([]int, *targetY+1+*margin)
		for y := 0; y <= *targetY+*margin; y++ {
			if (x == 0 && y == 0) || (x == *targetX && y == *targetY) {
				geoIdx[x][y] = 0
			} else if x == 0 {
				geoIdx[x][y] = y * 48271
			} else if y == 0 {
				geoIdx[x][y] = x * 16807
			} else {
				// This should already be memoized.
				geoIdx[x][y] = erosion[x-1][y] * erosion[x][y-1]
			}
			erosion[x][y] = (geoIdx[x][y] + *depth) % 20183
			terrain := erosion[x][y] % 3
			if x <= *targetX && y <= *targetY {
				risk += terrain
			}
		}
	}
	fmt.Printf("Total risk of area is %d\n", risk)

	start := Position{0, 0, false, true}
	winner := Position{*targetX, *targetY, false, true}

	shortestPath := ComputeShortest(start, winner, erosion)
	fmt.Printf("Took %d moves to reach our friend.\n", shortestPath)
}

func ComputeShortest(start, winner Position, erosion [][]int) int {
	seen := make(map[Position]int)
	seen[start] = 0
	queue := []Position{start}

	// Perform a BFS with backtracking.
	for len(queue) != 0 {
		pos := queue[0]
		queue = queue[1:]

		if shortest, ok := seen[pos]; !ok {
			// We don't know how to get here yet.
			continue
		} else {
			moves := pos.Moves(erosion)
			if len(moves) == 0 {
				// Illegal position.
				continue
			}
			for neighbor, incremental := range moves {
				shortestPathToNeighbor, ok := seen[neighbor]
				newPathLength := shortest + incremental
				if !ok || newPathLength < shortestPathToNeighbor {
					// Record new paths, as well as shortened paths.
					// Enqueue them to be re-run.
					seen[neighbor] = shortest + incremental
					queue = append(queue, neighbor)
				}
			}
		}
	}
	return seen[winner]
}
