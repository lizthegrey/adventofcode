package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day06.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")
var maxTotalDistance = flag.Int("maxTotalDistance", 10000, "Total distance threshold for Part B.")

type Coord struct {
	X, Y int
}

func (c Coord) Grow() []Coord {
	return []Coord{
		{c.X + 1, c.Y},
		{c.X - 1, c.Y},
		{c.X, c.Y + 1},
		{c.X, c.Y - 1},
	}
}

func (c Coord) Distance(c2 Coord) int {
	total := 0
	xdiff := c.X - c2.X
	ydiff := c.Y - c2.Y
	if xdiff < 0 {
		total -= xdiff
	} else {
		total += xdiff
	}
	if ydiff < 0 {
		total -= ydiff
	} else {
		total += ydiff
	}
	return total
}

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	points := make([]*Coord, 0)

	r := bufio.NewReader(f)
	for {
		l, err := r.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]
		parts := strings.Split(l, ", ")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		points = append(points, &Coord{x, y})
	}
	if !*partB {
		PartA(points)
	} else {
		PartB(points)
	}
}

func TotalDistance(points []*Coord, c Coord) int {
	sum := 0
	for _, p := range points {
		sum += p.Distance(c)
	}
	return sum
}

func PartB(points []*Coord) {
	seen := make(map[Coord]int)
	worklist := make([]Coord, 0)
	for _, i := range points {
		worklist = append(worklist, *i)
	}
	for len(worklist) != 0 {
		nextWorklist := make([]Coord, 0)
		for _, p := range worklist {
			if seen[p] != 0 {
				continue
			}
			d := TotalDistance(points, p)
			if d >= *maxTotalDistance {
				continue
			}
			seen[p] = d
			nextWorklist = append(nextWorklist, p.Grow()...)
		}
		worklist = nextWorklist
	}
	fmt.Printf("Result is %d.\n", len(seen))
}

func PartA(points []*Coord) {
	// key: coordinate being claimed
	// value: pointer to source node.
	claimed := make(map[Coord]*Coord)
	maxDistance := 0
	for i, c := range points {
		for j, c2 := range points {
			if i >= j {
				continue
			}
			if dist := c.Distance(*c2); dist > maxDistance {
				maxDistance = dist
			}
		}
	}
	maxDistance = maxDistance/2 + 1
	for _, c := range points {
		claimed[*c] = c
	}
	conflicts := make(map[Coord]bool)
	for distance := 1; distance < maxDistance; distance++ {
		extended := make(map[Coord]*Coord)
		for k, v := range claimed {
			neighbors := k.Grow()
			for _, n := range neighbors {
				if conflicts[n] == true {
					continue
				}
				extConflict := extended[n] != nil && extended[n] != v
				claimConflict := claimed[n] != nil && claimed[n] != v
				if extConflict || claimConflict {
					conflicts[n] = true
				} else if claimed[n] == v {
					continue
				} else {
					extended[n] = v
				}
			}
		}
		for k, v := range extended {
			if conflicts[k] {
				continue
			}
			claimed[k] = v
		}
		fmt.Printf("Finished pass %d of %d.\n", distance, maxDistance)
	}

	disqualified := make(map[*Coord]bool)
	for k, v := range claimed {
		neighbors := k.Grow()
		for _, n := range neighbors {
			if conflicts[n] == true {
				continue
			}
			if claimed[n] != nil {
				continue
			} else {
				disqualified[v] = true
				// Disqualify -- infinite.
			}
		}
	}

	sums := make(map[*Coord]int)
	for _, v := range claimed {
		sums[v] += 1
	}
	highest := 0
	for k, v := range sums {
		if disqualified[k] {
			fmt.Printf("Disqualified %d,%d\n", k.X, k.Y)
		} else {
			fmt.Printf("Coord %d,%d has area %d\n", k.X, k.Y, v)
			if highest < v {
				highest = v
			}
		}
	}
	fmt.Printf("Result is %d\n", highest)
}
