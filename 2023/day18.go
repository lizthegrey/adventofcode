package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"slices"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day18.input", "Relative file path to use as input.")

type State int

const (
	Unknown = iota
	Trench
	Floor
	GroundLevel
)

type Direction int

const (
	Right = iota
	Down
	Left
	Up
)

type Instruction struct {
	Steps int
	Dir   Direction
}
type Instructions []Instruction

type Coord struct {
	X, Y int
}

func (c Coord) Step(instr Instruction) Coord {
	var xDelta, yDelta int
	switch instr.Dir {
	case Up:
		yDelta = 1
	case Right:
		xDelta = 1
	case Left:
		xDelta = -1
	case Down:
		yDelta = -1
	}
	return Coord{c.X + xDelta*instr.Steps, c.Y + yDelta*instr.Steps}
}

type AxisStops []int

type BoundingBox struct {
	MinX, MinY, MaxX, MaxY int
	Status                 State
}

func (b BoundingBox) Area() uint64 {
	return uint64(b.MaxY-b.MinY+1) * uint64(b.MaxX-b.MinX+1)
}

type BoxCorner map[Coord]*BoundingBox

type Bounds struct {
	Boxes   []*BoundingBox
	Corners BoxCorner
}

// Basic theory: instead of using individual tiles as adjacencies,
// construct a set of axis stops; each axis stop splits the X and Y coordinates.
// Then we have bounding boxes for each: some that are edges, some of which are fill
// and some of which are empty, and we need to flood fill between adjacent bounding boxes
// but we don't need to treat each individual item inside a bounding box as unique.
//
// X    X  X
// 1ffffeddc YYY
// 2xxxxyzzb
// 344445wwa YYY
// vrrrr6uu0
// vrrrr6uu0
// vrrrr6uu0
// tssss7889 YYY
func MakeBoxes(xStops, yStops AxisStops) Bounds {
	ret := Bounds{
		Boxes:   nil,
		Corners: make(BoxCorner),
	}
	prevX := xStops[0] - 1
	for _, x := range xStops {
		if prevX == x {
			continue
		} else if prevX > x {
			log.Fatalf("AxisStops must be sorted first.")
		}
		prevY := yStops[0] - 1
		for _, y := range yStops {
			if prevY == y {
				continue
			} else if prevY > y {
				log.Fatalf("AxisStops must be sorted first.")
			}
			if prevX < x-1 && prevY < y-1 {
				// Create a new bounding box for the space between, before creating the stop's 1-width box.
				ret.AddBox(prevX+1, x-1, prevY+1, y-1)
			}
			if prevX < x-1 {
				// Create a new bounding box for the space between, before creating the stop's 1-width box.
				ret.AddBox(prevX+1, x-1, y, y)
			}
			if prevY < y-1 {
				// Create a new bounding box for the space between, before creating the stop's 1-width box.
				ret.AddBox(x, x, prevY+1, y-1)
			}
			// Intersection boxes will always be 1 wide.
			ret.AddBox(x, x, y, y)
			prevY = y
		}
		prevX = x
	}
	return ret
}

func (b *Bounds) AddBox(minX, maxX, minY, maxY int) {
	box := &BoundingBox{
		MinX: minX,
		MaxX: maxX,
		MinY: minY,
		MaxY: maxY,
	}
	b.Corners[Coord{minX, minY}] = box
	b.Corners[Coord{minX, maxY}] = box
	b.Corners[Coord{maxX, minY}] = box
	b.Corners[Coord{maxX, maxY}] = box
	b.Boxes = append(b.Boxes, box)
}

func (b Bounds) Area() uint64 {
	var sum uint64
	for _, box := range b.Boxes {
		if box.Status == Trench || box.Status == Floor {
			sum += box.Area()
		}
	}
	return sum
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var partA, partB Instructions
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, " ")
		var instrA, instrB Instruction
		switch parts[0] {
		case "U":
			instrA.Dir = Up
		case "R":
			instrA.Dir = Right
		case "D":
			instrA.Dir = Down
		case "L":
			instrA.Dir = Left
		}
		distance, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatalf("Failed parsing %s: %v", parts[1], err)
		}
		instrA.Steps = distance
		partA = append(partA, instrA)

		hex := parts[2][2 : len(parts[2])-1]
		val, err := strconv.ParseUint(hex, 16, 0)
		if err != nil {
			log.Fatalf("Failed parsing %s: %v", parts[2], err)
		}
		instrB.Steps = int(val >> 4)
		instrB.Dir = Direction(val & 3)
		partB = append(partB, instrB)
	}
	fmt.Println(partA.Process())
	fmt.Println(partB.Process())
}

func (instrs Instructions) Process() uint64 {
	// First, create the bounding boxes.
	var loc Coord
	xStops := AxisStops{0}
	yStops := AxisStops{0}
	for _, instr := range instrs {
		loc = loc.Step(instr)
		xStops = append(xStops, loc.X)
		yStops = append(yStops, loc.Y)
	}
	slices.Sort(xStops)
	slices.Sort(yStops)
	bounds := MakeBoxes(xStops, yStops)

	loc = Coord{0, 0}
	for _, instr := range instrs {
		dst := loc.Step(instr)
		for loc != dst {
			box, ok := bounds.Corners[loc]
			if !ok {
				log.Fatalf("Failed to find box for corner %v", loc)
			}
			box.Status = Trench
			switch instr.Dir {
			case Up:
				loc.Y = box.MaxY + 1
			case Right:
				loc.X = box.MaxX + 1
			case Down:
				loc.Y = box.MinY - 1
			case Left:
				loc.X = box.MinX - 1
			}
		}
	}
	// Mark the final tile as visited.
	box, ok := bounds.Corners[loc]
	if !ok {
		log.Fatalf("Failed to find box for final step %v", loc)
	}
	box.Status = Trench

	maxY := yStops[len(yStops)-1]

	var q []*BoundingBox
	var seenEdge bool
	for _, x := range xStops {
		loc := Coord{x, maxY - 1}
		if box := bounds.Corners[loc]; seenEdge && box.Status != Trench {
			q = append(q, box)
			break
		} else if box.Status == Trench {
			seenEdge = true
		}
	}

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur.Status != Unknown {
			continue
		}
		cur.Status = Floor
		for _, dir := range []Direction{Up, Right, Down, Left} {
			var move Coord
			switch dir {
			case Up:
				move = Coord{cur.MaxX + 0, cur.MaxY + 1}
			case Right:
				move = Coord{cur.MaxX + 1, cur.MaxY + 0}
			case Down:
				move = Coord{cur.MinX - 0, cur.MinY - 1}
			case Left:
				move = Coord{cur.MinX - 1, cur.MinY - 0}
			}
			neigh := bounds.Corners[move]
			if neigh != nil && neigh.Status != Trench {
				q = append(q, neigh)
			}
		}
	}
	return bounds.Area()
}
