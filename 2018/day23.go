package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day23.input", "Relative file path to use as input.")

type Coord struct {
	X, Y, Z int
}

type Drone struct {
	Loc   Coord
	Power int
}

func (c Coord) WalkPartway(t Coord, d int) Coord {
	xDelta := t.X - c.X
	yDelta := t.Y - c.Y
	zDelta := t.Z - c.Z

	sum := 0
	if xDelta < 0 {
		sum -= xDelta
	} else {
		sum += xDelta
	}
	if yDelta < 0 {
		sum -= yDelta
	} else {
		sum += yDelta
	}
	if zDelta < 0 {
		sum -= zDelta
	} else {
		sum += zDelta
	}

	return Coord{c.X + d*xDelta/sum, c.Y + d*yDelta/sum, c.Z + d*zDelta/sum}
}

func (c Coord) Distance(t Coord) int {
	xDelta := c.X - t.X
	yDelta := c.Y - t.Y
	zDelta := c.Z - t.Z
	if xDelta < 0 {
		xDelta = -xDelta
	}
	if yDelta < 0 {
		yDelta = -yDelta
	}
	if zDelta < 0 {
		zDelta = -zDelta
	}
	return xDelta + yDelta + zDelta
}

func (d Drone) InRange(t Coord) bool {
	return d.Loc.Distance(t) <= d.Power
}

func (d Drone) DistanceToRange(t Coord) int {
	return d.Loc.Distance(t) - d.Power
}

type Candidate struct {
	Loc      Coord
	Distance int
}

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	drones := make([]Drone, 0)

	r := bufio.NewReader(f)
	for {
		l, err := r.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]

		parts := strings.Split(l[5:], ">, r=")
		power, _ := strconv.Atoi(parts[1])
		position := strings.Split(parts[0], ",")
		x, _ := strconv.Atoi(position[0])
		y, _ := strconv.Atoi(position[1])
		z, _ := strconv.Atoi(position[2])
		drones = append(drones, Drone{Coord{x, y, z}, power})
	}

	highestPower := 0
	dronesInRangeForHighest := 0

	seeds := make(map[Coord]int)
	for _, d := range drones {
		dronesInRange := 0
		rangeToMe := 0
		for _, t := range drones {
			if d.InRange(t.Loc) {
				dronesInRange++
			}
			if t.InRange(d.Loc) {
				rangeToMe++
			}
		}
		seeds[d.Loc] = rangeToMe

		if d.Power > highestPower {
			highestPower = d.Power
			dronesInRangeForHighest = dronesInRange
		}
	}
	fmt.Printf("Highest drones in range of one drone: %d\n", dronesInRangeForHighest)

	ranges := make(map[Coord]int)
	worklist := make([]Candidate, 0)

	// Binary search the space, starting from the list of drone locations,
	// then walking to be in range of each of the N closest drones and adding to worklist.
	for l, n := range seeds {
		worklist = append(worklist, Candidate{l, n})
	}

	highScore := 0
	for len(worklist) > 0 {
		sort.Slice(worklist, func(i, j int) bool {
			return worklist[i].Distance > worklist[j].Distance
		})

		c := worklist[0]
		loc := c.Loc
		worklist = worklist[1:]

		if ranges[loc] != 0 {
			// Don't re-process items.
			continue
		}

		alreadyIn := 0
		candidates := make([]Candidate, 0)

		for _, d := range drones {
			if d.InRange(loc) {
				// We don't need to get any closer.
				alreadyIn++
				ranges[loc]++
				continue
			}
			dist := d.DistanceToRange(loc)
			mid := loc.WalkPartway(d.Loc, dist+1)
			candidates = append(candidates, Candidate{mid, dist})
		}

		if ranges[loc] < highScore - 2 {
			// This isn't worth bothering with.
			continue
		} else if ranges[loc] > highScore {
			fmt.Printf("New high score at %v: in range of %d\n", loc, ranges[loc])
			highScore = ranges[loc]
		}

		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Distance < candidates[j].Distance
		})

		for _, c := range candidates {
			worklist = append(worklist, Candidate{c.Loc, ranges[loc] + 1})
		}
	}

	lowestSum := 9999999999999
	var lowestCoord Coord
	for k, v := range ranges {
		if v != highScore {
			continue
		}
		mySum := k.X + k.Y + k.Z
		if mySum < lowestSum {
			lowestSum = mySum
			lowestCoord = k
		}
	}

	fmt.Printf("Checking if we can get closer on X Y Z: ")
	for _, v := range []*int{&lowestCoord.X, &lowestCoord.Y, &lowestCoord.Z} {
		original := *v
		for diff := 0; ; diff++ {
			*v--
			inRange := 0
			for _, d := range drones {
				if d.InRange(lowestCoord) {
					inRange++
				}
			}
			if inRange != highScore {
			fmt.Printf("%d ", diff)
				*v = original
				break
			}
		}
	}
	fmt.Println()

	fmt.Println(lowestCoord.X + lowestCoord.Y + lowestCoord.Z)
}
