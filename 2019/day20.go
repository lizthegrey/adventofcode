package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day20.input", "Relative file path to use as input.")

type State struct {
	Player      Coord
	MovesToDate int
}

type Coord struct {
	X, Y int
}

func (c Coord) Move(dir int, portals map[Coord]Coord) Coord {
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
	if jump, ok := portals[ret]; ok {
		return jump
	}
	return ret
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	passable := make(map[Coord]bool)
	portals := make(map[Coord]Coord)
	letters := make(map[Coord]rune)
	for y, s := range split {
		if s == "" {
			continue
		}
		for x, c := range s {
			loc := Coord{x, y}
			switch c {
			case '.':
				passable[loc] = true
			case ' ':
				// Do nothing.
			case '#':
				// Do nothing.
			default:
				letters[loc] = c
			}
		}
	}

	// Maps a pair of letters (e.g. AB) to two coordinates:
	// (0) the point we spit people out to
	// (1) the virtual square that you must move onto to be spit out at (0).
	entrances := make(map[[2]rune][2]Coord)
	for l, c := range letters {
		// If there is a letter above us, then it is the first char and we are the last.
		// Ditto if there is a letter to our left.
		// In that case, we'll find it on the other iteration.

		// Check whether we have a letter next to us (URDL)
		if _, ok := letters[Coord{l.X, l.Y - 1}]; ok {
			// Check up.
			continue
		}
		if _, ok := letters[Coord{l.X - 1, l.Y}]; ok {
			// Check left.
			continue
		}

		// If so, check which side is a passable . character.
		// Then, create a portal.
		var pair [2]rune
		var destination, source Coord
		pair[0] = c
		if o, ok := letters[Coord{l.X, l.Y + 1}]; ok {
			// Use the character below as the second of the identifier.
			pair[1] = o
			// Then check whether the portal end is below or above.
			downSpace := Coord{l.X, l.Y + 2}
			upSpace := Coord{l.X, l.Y - 1}
			if passable[downSpace] {
				source = Coord{l.X, l.Y + 1}
				destination = downSpace
			} else if passable[upSpace] {
				source = l
				destination = upSpace
			} else {
				fmt.Printf("Failed to parse grid correctly at %d,%d\n", l.X, l.Y)
			}
		} else if o, ok := letters[Coord{l.X + 1, l.Y}]; ok {
			// Use the character at right as the second of the identifier.
			pair[1] = o
			// Then check whether the portal end is right or left.
			rightSpace := Coord{l.X + 2, l.Y}
			leftSpace := Coord{l.X - 1, l.Y}
			if passable[rightSpace] {
				source = Coord{l.X + 1, l.Y}
				destination = rightSpace
			} else if passable[leftSpace] {
				source = l
				destination = leftSpace
			} else {
				fmt.Printf("Failed to parse grid correctly at %d,%d\n", l.X, l.Y)
			}
		} else {
			fmt.Printf("Failed to parse grid correctly at %d,%d\n", l.X, l.Y)
		}

		// Create a dangling pair, or connect the pairs.
		if opposite, ok := entrances[pair]; ok {
			portals[source] = opposite[0]
			portals[opposite[1]] = destination
			delete(entrances, pair)
		} else {
			entrances[pair] = [2]Coord{destination, source}
		}
	}

	if len(entrances) != 2 {
		fmt.Printf("Failed to pair at least some portals: %v\n", portals)
	}
	start := entrances[[2]rune{'A', 'A'}][0]
	end := entrances[[2]rune{'Z', 'Z'}][0]

	shortest := bfs(start, end, passable, portals)
	fmt.Printf("Shortest path is %d long.\n", shortest)
}

type AStarItem struct {
	C    Coord
	Cost int
}

func (c Coord) Distance(o Coord) int {
	sum := 0
	if c.X < o.X {
		sum += o.X - c.X
	} else {
		sum += c.X - o.X
	}

	if c.Y < o.Y {
		sum += o.Y - c.Y
	} else {
		sum += c.Y - o.Y
	}
	return sum
}

func minDistanceToPortalOrEnd(src, dst Coord, portals map[Coord]Coord) int {
	minDistance := src.Distance(dst)
	for in := range portals {
		// TODO: can be improved by taking into account closeness of portal exits to
		// other portal entrances and/or to the finish.
		d := src.Distance(in) + 1
		if d < minDistance {
			minDistance = d
		}
	}
	return minDistance
}

func bfs(src, dst Coord, passable map[Coord]bool, portals map[Coord]Coord) int {
	// Perform a breadth-first search.
	shortest := map[Coord]int{
		src: 0,
	}
	worklist := []AStarItem{{src, minDistanceToPortalOrEnd(src, dst, portals)}}
	for {
		sort.Slice(worklist, func(i, j int) bool {
			return worklist[i].Cost < worklist[j].Cost
		})
		w := worklist[0].C
		for dir := 1; dir <= 4; dir++ {
			moved := w.Move(dir, portals)
			if !passable[moved] {
				// Don't check impassable tiles.
				continue
			}
			if _, ok := shortest[moved]; ok {
				continue
			} else {
				shortest[moved] = shortest[w] + 1
				if moved == dst {
					return shortest[moved]
				}
				worklist = append(worklist, AStarItem{moved, shortest[moved] + minDistanceToPortalOrEnd(moved, dst, portals)})
			}
		}
		if len(worklist) == 1 {
			// Not reachable; let the caller know.
			return -1
		}
		worklist = worklist[1:]
	}
}
