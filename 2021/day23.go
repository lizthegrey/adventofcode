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
var debug = flag.Bool("debug", false, "Whether to print debug output.")

var tr = otel.Tracer("day23")

type Terrain map[Coord]bool

type AmphipodType uint8

const (
	Undef AmphipodType = iota
	A
	B
	C
	D
)

type Amphipod struct {
	Loc  Coord
	Type AmphipodType
}

func (a Amphipod) MovementCost() int {
	switch a.Type {
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

func (a Amphipod) DestinationColumn() int8 {
	switch a.Type {
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
	R, C int8
}

type World struct {
	Pieces   [16]Amphipod
	HallPass int8
}

func (w World) PrintBoard(passable Terrain) {
	for r := 0; r < 7; r++ {
		for c := 0; c < 13; c++ {
			loc := Coord{int8(r), int8(c)}
			valid, contents := w.PassableOrContents(loc, passable)
			if contents != nil {
				fmt.Printf("%s", []byte{byte('A') + byte(contents.Type-1)})
			} else if valid {
				fmt.Printf(".")
			} else {
				fmt.Printf("#")
			}
		}
		fmt.Println()
	}
}

func (w World) GetPieces() []Amphipod {
	pieces := make([]Amphipod, 0, 16)
	for _, v := range w.Pieces {
		if v.Type == Undef {
			continue
		}
		pieces = append(pieces, v)
	}
	return pieces
}

func (w World) RouteClear(whoami Amphipod, passable Terrain) bool {
	r := whoami.Loc.R
	c := whoami.Loc.C
	if r != 1 {
		fmt.Println("We don't need a hall pass to move outside hallway.")
		return true
	}
	desiredColumn := whoami.DestinationColumn()
	// Try moving in a straight line, without moving any other pieces.
	// First, try moving left or right until we get to the column we want.
	var delta int8
	if desiredColumn > c {
		delta = int8(1)
	} else if desiredColumn < c {
		delta = int8(-1)
	}
	for col := c; col != desiredColumn; col += delta {
		if col == c {
			// Ignore ourself.
			continue
		}
		if !w.Passable(Coord{r, col}, passable) {
			// Something's in the way!
			return false
		}
	}
	// Verify there are no mismatched pieces in our destination.
	row := int8(2)
	for {
		empty, contents := w.PassableOrContents(Coord{row, desiredColumn}, passable)
		if !empty && contents == nil {
			// We've hit the bottom and can break out to return our true.
			break
		}
		if contents != nil && contents.Type != whoami.Type {
			return false
		}
		row++
	}
	return true
}

func (w World) HasReachedDestination(whoami Amphipod, passable Terrain) bool {
	r := whoami.Loc.R
	c := whoami.Loc.C
	if c != whoami.DestinationColumn() || r == 1 {
		// Wrong column, or still in the hallway.
		return false
	}
	row := r
	for {
		row++
		empty, contents := w.PassableOrContents(Coord{row, c}, passable)
		if !empty && contents == nil {
			// We're packed to the very bottom, so we're good.
			break
		}
		if empty {
			// There is empty space below us, keep moving!
			return false
		}
		// A slot below us is full, check whether it matches us.
		if contents.Type != whoami.Type {
			// We're going to need to move up to let someone out.
			return false
		}
	}
	// Yup, we're the last one in!
	return true
}

func (w World) Passable(coord Coord, passable Terrain) bool {
	ret, _ := w.PassableOrContents(coord, passable)
	return ret
}

func (w World) PassableOrContents(coord Coord, passable Terrain) (bool, *Amphipod) {
	if !passable[coord] {
		return false, nil
	}
	for _, v := range w.GetPieces() {
		if v.Loc == coord {
			return false, &v
		}
	}
	return true, nil
}

type FastMove struct {
	A Amphipod
	C int
}

func (w World) CanStopHere(src, proposed Amphipod, passable Terrain) bool {
	r := proposed.Loc.R
	c := proposed.Loc.C
	if r == 1 {
		if c == 3 || c == 5 || c == 7 || c == 9 {
			// It is not valid to stop in the doorway.
			return false
		} else {
			// It is valid to stop in between doorways or in an alcove, if and only
			// if I did not start on row 1; even if I hold a hall pass, once I stop
			// I can't move within the hallway, only go to my burrow.
			return src.Loc.R != 1
		}
	} else {
		// We're in a burrow. In general, we can't stop here once we've started
		// moving, unless this is at the very bottom of our terminal burrow.
		return w.HasReachedDestination(proposed, passable)
	}
}

// FastAllowedMovements fast-forwards a string of movements to a valid landing
// position, rather than returning one step at a time. We recursively use
// AllowedMovements to move ourselves, and place any landing position we find
// into our return list.
func (w World) FastAllowedMovements(src Amphipod, passable Terrain) ([]Amphipod, []int) {
	var ret []Amphipod
	var costs []int

	open := []FastMove{{src, 0}}
	visited := make(map[Amphipod]bool)
	for len(open) != 0 {
		popped := open[0]
		open = open[1:]

		if visited[popped.A] {
			continue
		}
		visited[popped.A] = true

		if src != popped.A && w.CanStopHere(src, popped.A, passable) {
			ret = append(ret, popped.A)
			costs = append(costs, popped.C)
			// But we still have to search neighbors of ourself just to be safe.
		}

		for _, dst := range w.AllowedMovements(popped.A, passable) {
			open = append(open, FastMove{dst, popped.C + popped.A.MovementCost()})
		}
	}
	return ret, costs
}

func (w World) AllowedMovements(a Amphipod, passable Terrain) []Amphipod {
	// Pieces can move up, down, left, right.
	// Pieces also cannot stop above a column, but we'll deal with
	// that elsewhere.
	var ret []Amphipod
	r := a.Loc.R
	c := a.Loc.C
	if dst := (Coord{r - 1, c}); w.Passable(dst, passable) {
		ret = append(ret, Amphipod{dst, a.Type})
	}
	if dst := (Coord{r + 1, c}); w.Passable(dst, passable) {
		// However, we cannot move into a column it does not belong in,
		// nor block the path for an Amphipod that needs to leave.
		valid := c == a.DestinationColumn()
		_, twoBelow := w.PassableOrContents(Coord{r + 2, c}, passable)
		if twoBelow != nil && twoBelow.DestinationColumn() != c {
			valid = false
		}
		if valid {
			ret = append(ret, Amphipod{dst, a.Type})
		}
	}
	if dst := (Coord{r, c - 1}); w.Passable(dst, passable) {
		ret = append(ret, Amphipod{dst, a.Type})
	}
	if dst := (Coord{r, c + 1}); w.Passable(dst, passable) {
		ret = append(ret, Amphipod{dst, a.Type})
	}
	return ret
}

func (w World) PieceWithForcedMovement() *Amphipod {
	for _, p := range w.GetPieces() {
		r := p.Loc.R
		c := p.Loc.C
		if r != 1 {
			continue
		}
		if c == 3 || c == 5 || c == 7 || c == 9 {
			return &p
		}
	}
	return nil
}

func (w World) GenerateMoves(passable Terrain) ([]World, []int) {
	var costs []int
	var moves []World

	forced := w.PieceWithForcedMovement()
	for i, src := range w.GetPieces() {
		if forced != nil && *forced != src {
			// We have to keep a piece moving off the entrance.
			continue
		}

		if w.HasReachedDestination(src, passable) {
			// We don't need to move, we're already in the right place.
			continue
		}

		// Only allow changing the HallPass holder if the route is clear to a
		// target space for the potential new HallPass holder.
		if w.HallPass != int8(-1) && w.HallPass != int8(i) {
			if src.Loc.R == 1 && !w.RouteClear(src, passable) {
				// We're stuck in place in the hallway for now, because there's no
				// route to our burrow.
				continue
			}
		}

		dsts, dstCosts := w.FastAllowedMovements(src, passable)
		for j, dst := range dsts {
			// Remove the tile from src, add it to dst.
			after := w
			after.Pieces[i] = dst
			// We need to null the HallPass holder (to -1) if the last
			// HallPass holder is no longer in the hall.
			if after.HasReachedDestination(dst, passable) {
				after.HallPass = int8(-1)
			}
			if dst.Loc.R == 1 {
				// We've moved in the hallway, so lock everyone else out.
				after.HallPass = int8(i)
			}
			moves = append(moves, after)
			costs = append(costs, dstCosts[j])
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
	for _, p := range w.GetPieces() {
		r := p.Loc.R
		c := p.Loc.C
		columnDiff := int(c - p.DestinationColumn())
		if columnDiff != 0 {
			// We need to move it up, over, and down.
			if columnDiff < 0 {
				columnDiff = -columnDiff
			}
			cost += p.MovementCost() * (1 + (int(r) - 1) + columnDiff)
		} else if r == 1 {
			// the cost is just the cost to slot it down into the desired location.
			cost += p.MovementCost()
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
	passable := make(Terrain)

	piecesFound := 0
	for r, line := range split {
		for c, v := range line {
			loc := Coord{int8(r), int8(c)}
			var kind AmphipodType
			switch v {
			case '#':
				fallthrough
			case ' ':
				// Not passable, don't put anything into passable map.
				continue
			case '.':
				passable[loc] = true
				continue
			case 'A':
				kind = A
			case 'B':
				kind = B
			case 'C':
				kind = C
			case 'D':
				kind = D
			default:
				fmt.Println("Encountered unknown board character.")
			}
			initial.Pieces[piecesFound] = Amphipod{loc, kind}
			passable[loc] = true
			piecesFound++
		}
	}

	f, err := os.Create("/tmp/pprof")
	defer f.Close() // error handling omitted for example
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Printf("could not start CPU profile: %v", err)
		return
	}
	defer pprof.StopCPUProfile()
	fmt.Println(AStar(initial, passable))

	// Manipulate board for part b.
	// Add additional passable tiles.
	for r := 4; r <= 5; r++ {
		for c := 3; c <= 9; c += 2 {
			passable[Coord{int8(r), int8(c)}] = true
		}
	}
	// Move down amphipods 4-7 by 2 rows.
	for i := range initial.Pieces {
		if i < 4 || i >= 8 {
			continue
		}
		initial.Pieces[i].Loc.R += 2
	}
	// And also insert additional rows.
	//   #D#C#B#A#
	//   #D#B#A#C#
	idx := 8
	for r, row := range [2][4]AmphipodType{{D, C, B, A}, {D, B, A, C}} {
		for c, kind := range row {
			initial.Pieces[idx] = Amphipod{Coord{int8(3 + r), int8(3 + 2*c)}, kind}
			idx++
		}
	}
	if *debug {
		fmt.Println("\n-------")
	}
	fmt.Println(AStar(initial, passable))
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

func AStar(src World, passable Terrain) int {
	gScore := CostMap{
		src: 0,
	}
	fScore := CostMap{
		src: src.MinCostToSort(),
	}
	from := make(map[World]*World)
	workList := HeapQueue{&[]World{src}, fScore, make(CostMap)}
	heap.Init(&workList)

	for len(*workList.Elems) != 0 {
		// Pop the current node off the worklist.
		current := heap.Pop(&workList).(World)

		if current.MinCostToSort() == 0 {
			if *debug {
				totalCost := gScore[current]
				prev := *from[current]
				current.PrintBoard(passable)
				for prev != src {
					fmt.Printf("Move cost %d (%d cumulative) to get from\n\n", totalCost-gScore[prev], totalCost)
					prev.PrintBoard(passable)
					totalCost = gScore[prev]
					prev = *from[prev]
				}
				fmt.Printf("Move cost %d (%d cumulative) to get from\n\n", totalCost, totalCost)
				prev.PrintBoard(passable)
			}
			return gScore[current]
		}
		moves, costs := current.GenerateMoves(passable)
		for i, after := range moves {
			proposedScore := gScore[current] + costs[i]
			if previousScore, ok := gScore[after]; !ok || proposedScore < previousScore {
				gScore[after] = proposedScore
				fScore[after] = proposedScore + after.MinCostToSort()
				from[after] = &current
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
