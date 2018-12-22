package main

import (
	"flag"
	"fmt"
)

var depth = flag.Int("depth", 8112, "The depth of the cave system.")
var targetX = flag.Int("targetX", 13, "The target's X coordinate.")
var targetY = flag.Int("targetY", 743, "The target's Y coordinate.")
var margin = flag.Int("margin", 150, "The margin by which we'll check 'longer' paths.")

type Move struct {
	X, Y  int
	Climb bool
	Torch bool
}

type Path struct {
	Steps    []Move
	Distance int
}

func (m Move) Passable(erosion [][]int) bool {
	if m.X < 0 || m.X >= len(erosion) {
		return false
	}
	if m.Y < 0 || m.Y >= len(erosion[0]) {
		return false
	}

	terrain := erosion[m.X][m.Y] % 3
	switch terrain {
	case 0:
		// Rocky
		return m.Climb || m.Torch
	case 1:
		// Wet
		return !m.Torch
	case 2:
		// Narrow
		return !m.Climb
	}
	return false
}

func (m Move) Moves(erosion [][]int) map[Move]int {
	ret := make(map[Move]int)
	terrain := erosion[m.X][m.Y] % 3
	swapGear := m
	switch terrain {
	case 0:
		// Rocky
		if !m.Climb && !m.Torch {
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

	up := m
	up.Y--
	if up.Passable(erosion) {
		ret[up] = 1
	}

	down := m
	down.Y++
	if down.Passable(erosion) {
		ret[down] = 1
	}

	left := m
	left.X--
	if left.Passable(erosion) {
		ret[left] = 1
	}

	right := m
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

	winner := Move{*targetX, *targetY, false, true}
	start := Move{0, 0, false, true}

	seen := make(map[Move]int)
	seen[start] = 0
	shortestPath := ComputeShortest(start, winner, erosion, seen)
	fmt.Printf("Took %d moves to reach our friend.\n", shortestPath)
}

func ComputeShortest(start, winner Move, erosion [][]int, seen map[Move]int) int {
	for maxSquareSide := 0; maxSquareSide < len(erosion) || maxSquareSide < len(erosion[0]); maxSquareSide++ {
		if maxSquareSide%10 == 0 {
			fmt.Printf("Processing square %d\n", maxSquareSide)
		}

		maxX := len(erosion)
		maxY := len(erosion[0])
		if maxSquareSide < maxX {
			maxX = maxSquareSide
		}
		if maxSquareSide < maxY {
			maxY = maxSquareSide
		}

		// Progressively re-compute the shortest paths to each location.
		queue := make([]Move, 0)
		for x := 0; x <= maxX; x++ {
			nothing := Move{x, maxY, false, false}
			climb := Move{x, maxY, true, false}
			torch := Move{x, maxY, false, true}
			queue = append(queue, nothing, climb, torch)
		}
		for y := 0; y <= maxY; y++ {
			nothing := Move{maxX, y, false, false}
			climb := Move{maxX, y, true, false}
			torch := Move{maxX, y, false, true}
			queue = append(queue, nothing, climb, torch)
		}

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
						seen[neighbor] = shortest + incremental
						queue = append(queue, neighbor)
					}
				}
			}
		}
	}
	return seen[winner]
}
