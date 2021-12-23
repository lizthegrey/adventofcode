package main

import (
	"container/heap"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/lizthegrey/adventofcode/2021/trace"
	"go.opentelemetry.io/otel"
	//"go.opentelemetry.io/otel/attribute"

	"os"
	"runtime/pprof"
)

var inputFile = flag.String("inputFile", "inputs/day23.input", "Relative file path to use as input.")

var tr = otel.Tracer("day23")

type Contents int

const (
	Undefined Contents = iota
	Empty
	Wall
	A
	B
	C
	D
)

func (c Contents) MovementCost() int {
	switch c {
	case A:
		return 1
	case B:
		return 10
	case C:
		return 100
	case D:
		return 1000
	default:
		fmt.Println("Asked for movement cost of a non-moveable token.")
		return -1
	}
}

func (c Contents) IsToken() bool {
	return c >= A && c <= D
}

func (c Contents) DestinationColumn() int {
	switch c {
	case A:
		return 3
	case B:
		return 5
	case C:
		return 7
	case D:
		return 9
	default:
		fmt.Println("Asked for destination column of a non-moveable token.")
		return -1
	}
}

type Coord struct {
	R, C int
}

type World struct {
	Tiles    [5][13]Contents
	Pieces   [8]Coord
	HallPass int
}

func (w World) RouteClear(coord Coord) bool {
	r := coord.R
	c := coord.C
	if r != 1 {
		fmt.Println("We don't need a hall pass to move outside hallway.")
		return true
	}
	whoami := w.Tiles[r][c]
	desiredColumn := whoami.DestinationColumn()
	// Try moving in a straight line, without moving any other pieces.
	// First, try moving left or right until we get to the column we want.
	var delta int
	if desiredColumn > c {
		delta = 1
	} else if desiredColumn < c {
		delta = -1
	}
	for col := c; col != desiredColumn; col += delta {
		if col == c {
			// Ignore ourself.
			continue
		}
		if w.Tiles[r][col] != Empty {
			// Something's in the way!
			return false
		}
	}
	// Verify there are no mismatched pieces in our destination.
	for row := 2; row <= 3; row++ {
		if slot := w.Tiles[2][desiredColumn]; slot != Empty && slot != whoami {
			return false
		}
	}
	return true
}

func (w World) HasReachedDestination(coord Coord) bool {
	r := coord.R
	c := coord.C
	whoami := w.Tiles[r][c]
	if c != whoami.DestinationColumn() || r == 1 {
		// Wrong column, or still in the hallway.
		return false
	}
	if r == 3 {
		// We've pushed all the way down, and we're in the right place.
		return true
	}
	if w.Tiles[3][c] == Empty {
		// There is empty space below us, keep moving!
		return false
	}
	// The slot below us is full
	if w.Tiles[3][c] != whoami {
		// We're going to need to move up to let someone out.
		return false
	}
	// Yup, we're the last one in!
	return true
}

func (w World) AllowedMovements(coord Coord) []Coord {
	// Pieces can move up, down, left, right.
	// Pieces also cannot stop above a column, but we'll deal with
	// that elsewhere.
	var ret []Coord
	r := coord.R
	c := coord.C
	if w.Tiles[r-1][c] == Empty {
		ret = append(ret, Coord{r - 1, c})
	}
	if w.Tiles[r+1][c] == Empty {
		// However, we cannot move into a column it does not belong in,
		// nor block the path for an Amphipod that needs to leave.
		valid := c == w.Tiles[r][c].DestinationColumn()
		twoBelow := w.Tiles[r+2][c]
		if twoBelow.IsToken() && twoBelow.DestinationColumn() != c {
			valid = false
		}
		if valid {
			ret = append(ret, Coord{r + 1, c})
		}
	}
	if w.Tiles[r][c-1] == Empty {
		ret = append(ret, Coord{r, c - 1})
	}
	if w.Tiles[r][c+1] == Empty {
		ret = append(ret, Coord{r, c + 1})
	}
	return ret
}

func (w World) PieceWithForcedMovement() *Coord {
	r := 1
	for c := 3; c <= 9; c += 2 {
		if w.Tiles[r][c] != Empty {
			return &Coord{r, c}
		}
	}
	return nil
}

func (w World) GenerateMoves() ([]World, []int) {
	var costs []int
	var moves []World

	forced := w.PieceWithForcedMovement()
	for i, src := range w.Pieces {
		if forced != nil && *forced != src {
			// We have to keep a piece moving off the entrance.
			continue
		}

		if w.HasReachedDestination(src) {
			// We don't need to move, we're already in the right place.
			continue
		}

		// Only allow changing the HallPass holder if the route is clear to a
		// target space for the potential new HallPass holder.
		if w.HallPass != -1 && w.HallPass != i {
			if src.R == 1 && !w.RouteClear(src) {
				// We're stuck in place in the hallway for now, because there's no
				// route to our burrow.
				continue
			}
		}

		kind := w.Tiles[src.R][src.C]
		for _, dst := range w.AllowedMovements(src) {
			// Remove the tile from src, add it to dst.
			after := w
			after.Tiles[src.R][src.C] = Empty
			after.Tiles[dst.R][dst.C] = kind
			after.Pieces[i] = dst
			// We need to null the HallPass holder (to -1) if the last
			// HallPass holder is no longer in the hall.
			if after.HasReachedDestination(dst) {
				after.HallPass = -1
			}
			if dst.R == 1 {
				// We've moved in the hallway, so lock everyone else out.
				after.HallPass = i
			}
			moves = append(moves, after)
			costs = append(costs, kind.MovementCost())
		}
	}

	return moves, costs
}

// MinCostToSort returns 0 if already sorted, otherwise provides a heuristic
// for the minimum cost to sort each piece. This assumes pieces can move through
// each other, but not through walls, and will only have to land in the top
// available slot.
func (w World) MinCostToSort() int {
	var cost int
	for _, p := range w.Pieces {
		kind := w.Tiles[p.R][p.C]
		columnDiff := p.C - kind.DestinationColumn()
		if columnDiff != 0 {
			// We need to move it up, over, and down.
			if columnDiff < 0 {
				columnDiff = -columnDiff
			}
			cost += kind.MovementCost() * (1 + (p.R - 1) + columnDiff)
		} else if p.R == 1 {
			// the cost is just the cost to slot it down into the desired location.
			cost += kind.MovementCost()
		} else {
			// We're already in the desired location, 0 cost.
		}
	}
	return cost
}

func main() {
	flag.Parse()

	ctx := context.Background()
	hny, tp := trace.InitializeTracing(ctx)
	defer hny.Shutdown(ctx)
	defer tp.Shutdown(ctx)

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	var initial World
	initial.HallPass = -1

	piecesFound := 0
	for r, line := range split {
		for c, v := range line {
			switch v {
			case '#':
				initial.Tiles[r][c] = Wall
			case 'A':
				initial.Tiles[r][c] = A
				initial.Pieces[piecesFound] = Coord{r, c}
				piecesFound++
			case 'B':
				initial.Tiles[r][c] = B
				initial.Pieces[piecesFound] = Coord{r, c}
				piecesFound++
			case 'C':
				initial.Tiles[r][c] = C
				initial.Pieces[piecesFound] = Coord{r, c}
				piecesFound++
			case 'D':
				initial.Tiles[r][c] = D
				initial.Pieces[piecesFound] = Coord{r, c}
				piecesFound++
			case '.':
				initial.Tiles[r][c] = Empty
			case ' ':
				// Skip
			default:
				fmt.Println("Encountered unknown board character.")
			}
		}
	}

	f, err := os.Create("/tmp/pprof")
	defer f.Close() // error handling omitted for example
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Printf("could not start CPU profile: %v", err)
		return
	}
	defer pprof.StopCPUProfile()
	fmt.Println(AStar(initial))
}

type CostMap map[World]int
type HeapQueue struct {
	Elems            *[]World
	Score, Positions CostMap
}

func (h HeapQueue) Len() int           { return len(*h.Elems) }
func (h HeapQueue) Less(i, j int) bool { return h.Score[(*h.Elems)[i]] < h.Score[(*h.Elems)[j]] }
func (h HeapQueue) Swap(i, j int) {
	h.Positions[(*h.Elems)[i]], h.Positions[(*h.Elems)[j]] = h.Positions[(*h.Elems)[j]], h.Positions[(*h.Elems)[i]]
	(*h.Elems)[i], (*h.Elems)[j] = (*h.Elems)[j], (*h.Elems)[i]
}

func (h HeapQueue) Push(x interface{}) {
	h.Positions[x.(World)] = len(*h.Elems)
	*h.Elems = append(*h.Elems, x.(World))
}

func (h HeapQueue) Pop() interface{} {
	old := *h.Elems
	n := len(old)
	x := old[n-1]
	*h.Elems = old[0 : n-1]
	delete(h.Positions, x)
	return x
}

func (h HeapQueue) Position(x World) int {
	if pos, ok := h.Positions[x]; ok {
		return pos
	}
	return -1
}

func AStar(src World) int {
	gScore := CostMap{
		src: 0,
	}
	fScore := CostMap{
		src: src.MinCostToSort(),
	}
	workList := HeapQueue{&[]World{src}, fScore, make(CostMap)}
	heap.Init(&workList)

	for len(*workList.Elems) != 0 {
		// Pop the current node off the worklist.
		current := heap.Pop(&workList).(World)

		if current.MinCostToSort() == 0 {
			return gScore[current]
		}
		moves, costs := current.GenerateMoves()
		for i, after := range moves {
			proposedScore := gScore[current] + costs[i]
			if previousScore, ok := gScore[after]; !ok || proposedScore < previousScore {
				gScore[after] = proposedScore
				fScore[after] = proposedScore + after.MinCostToSort()
				if pos := workList.Position(after); pos == -1 {
					heap.Push(&workList, after)
				} else {
					heap.Fix(&workList, pos)
				}
			}
		}
	}
	return -1
}
