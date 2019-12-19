package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day18.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

type CacheKey struct {
	Player        [4]Coord
	KeysCollected [26]bool
}

type State struct {
	Player      [4]Coord
	MovesToDate int
	// Note: this is in order.
	KeyCount             int
	KeysCollectedInOrder [26]int
	KeysCollected        [26]bool
}

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
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var center Coord
	var player [4]Coord
	passable := make(map[Coord]bool)
	keys := make(map[Coord]int)
	var keyList [26]Coord
	doors := make(map[Coord]int)
	for y, s := range split {
		if s == "" {
			continue
		}
		for x, c := range s {
			loc := Coord{x, y}
			switch c {
			case '#':
				passable[loc] = false
			case '@':
				center = loc
				fallthrough
			case '.':
				passable[loc] = true
			default:
				if c >= 'a' {
					keys[loc] = int(c - 'a')
					keyList[int(c-'a')] = loc
				} else {
					doors[loc] = int(c - 'A')
				}
				passable[loc] = true
			}
		}
	}

	if !*partB {
		player[0] = center
	} else {
		passable[center] = false
		passable[Coord{center.X + 1, center.Y}] = false
		passable[Coord{center.X - 1, center.Y}] = false
		passable[Coord{center.X, center.Y + 1}] = false
		passable[Coord{center.X, center.Y - 1}] = false
		player[0] = Coord{center.X + 1, center.Y + 1}
		player[1] = Coord{center.X - 1, center.Y + 1}
		player[2] = Coord{center.X - 1, center.Y - 1}
		player[3] = Coord{center.X + 1, center.Y - 1}
	}

	// Determine which keys are reachable based off current unlocks.
	// Go pick up one of the keys. (e.g. queue in priority order)
	// Determine which keys are reachable.
	// Once all keys are collected, report the number of moves taken.
	var emptyKeys [26]bool
	var emptyKeyList [26]int
	for i := range emptyKeyList {
		emptyKeyList[i] = -1
	}
	initial := State{player, 0, 0, emptyKeyList, emptyKeys}
	stateCache := make(map[CacheKey]int)
	unexplored := []State{initial}
	shortest := math.MaxInt32
	keyCount := -1
	for {
		state := unexplored[0]
		ck := CacheKey{state.Player, state.KeysCollected}

		if shortest, ok := stateCache[ck]; ok && shortest <= state.MovesToDate {
			if len(unexplored) == 1 {
				// Finished exploring the maze.
				break
			}
			unexplored = unexplored[1:]
			continue
		}
		stateCache[ck] = state.MovesToDate

		if state.KeyCount > keyCount {
			fmt.Printf("Progress: %d of %d\n", state.KeyCount, len(keys))
			keyCount = state.KeyCount
		}

		for i, found := range state.KeysCollected {
			if found {
				// Already collected this key.
				continue
			}
			if i >= len(keys) {
				continue
			}
			// calculate whether the key is reachable based on what's reachable so far.
			// treat the key we're not looking for as impassable to avoid accounting errors.
			quad := 0
			if *partB {
				quad = quadrant(center, keyList[i])
			}
			start := state.Player[quad]
			moves := bfs(start, keyList[i], passable, state.KeysCollected, keys, doors)
			if moves < 0 {
				// Not reachable.
				continue
			}

			newState := state
			newState.Player[quad] = keyList[i]
			newState.KeysCollectedInOrder[state.KeyCount] = i
			newState.KeysCollected[i] = true
			newState.MovesToDate += moves
			newState.KeyCount++

			if newState.MovesToDate >= shortest {
				// No point in searching futile states.
				continue
			}

			if newState.KeyCount == len(keys) {
				if newState.MovesToDate < shortest {
					shortest = newState.MovesToDate
				}
				continue
			}

			unexplored = append(unexplored, newState)
		}

		if len(unexplored) == 1 {
			// Finished exploring the maze.
			break
		}
		unexplored = unexplored[1:]
	}
	fmt.Printf("Shortest path is %d long.\n", shortest)
}

func quadrant(c, k Coord) int {
	// player[0] = Coord{center.X + 1, center.Y + 1}
	// player[1] = Coord{center.X - 1, center.Y + 1}
	// player[2] = Coord{center.X - 1, center.Y - 1}
	// player[3] = Coord{center.X + 1, center.Y - 1}
	if c.X < k.X {
		if c.Y < k.Y {
			return 0
		} else {
			return 3
		}
	} else {
		if c.Y < k.Y {
			return 1
		} else {
			return 2
		}
	}
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

func bfs(src, dst Coord, passable map[Coord]bool, keyState [26]bool, keys, doors map[Coord]int) int {
	// Perform a breadth-first search.
	shortest := map[Coord]int{
		src: 0,
	}
	worklist := []AStarItem{{src, src.Distance(dst)}}
	for {
		sort.Slice(worklist, func(i, j int) bool {
			return worklist[i].Cost < worklist[j].Cost
		})
		w := worklist[0].C
		for dir := 1; dir <= 4; dir++ {
			moved := w.Move(dir)
			if !passable[moved] {
				// Don't check impassable tiles.
				continue
			}
			if id, ok := keys[moved]; ok {
				// Don't allow moving through keys we haven't picked up yet.
				// but allow moving into the square that we're trying to get to.
				if !keyState[id] && moved != dst {
					continue
				}
			}
			if id, ok := doors[moved]; ok {
				// Don't allow moving through doors we don't have the key for.
				if !keyState[id] {
					continue
				}
			}
			if _, ok := shortest[moved]; ok {
				continue
			} else {
				shortest[moved] = shortest[w] + 1
				if moved == dst {
					return shortest[moved]
				}
				worklist = append(worklist, AStarItem{moved, shortest[moved] + moved.Distance(dst)})
			}
		}
		if len(worklist) == 1 {
			// Not reachable; let the caller know.
			return -1
		}
		worklist = worklist[1:]
	}
}
