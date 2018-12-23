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
	dronesInRange := 0
	for _, d := range drones {
		if d.Power > highestPower {
			highestPower = d.Power
			dronesInRange = 0
			for _, t := range drones {
				if d.InRange(t.Loc) {
					dronesInRange++
				}
			}
		}
	}
	fmt.Printf("Highest drones in range of one drone: %d\n", dronesInRange)

	ranges := make(map[Coord]int)
	worklist := make([]Candidate, 0)
	// Binary search the space, starting from the list of drone locations,
	// then creating a midpoint to the 10 closest drones and adding to worklist.
	for _, d := range drones {
		worklist = append(worklist, Candidate{d.Loc, 0})
	}

	highScore := 823
	//highScore := 0
	for len(worklist) > 0 {
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

		if ranges[loc] < highScore {
			// This isn't worth bothering with.
			continue
		} else {
			highScore = ranges[loc]
		}

		fmt.Printf("Checking location %v: in range of %d\n", loc, ranges[loc])

		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Distance < candidates[j].Distance
		})

		for i, c := range candidates {
			if i > 15 {
				// Don't fork too much initially.
				break
			}
			worklist = append(worklist, Candidate{c.Loc, ranges[loc] + 1})

		}
		sort.Slice(worklist, func(i, j int) bool {
			return worklist[i].Distance > worklist[j].Distance
		})
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

	for diff := 0; ; diff++ {
		loc := lowestCoord
		loc.X -= diff
		inRange := 0
		for _, d := range drones {
			if d.InRange(loc) {
				inRange++
			}
		}
		if inRange != highScore {
			fmt.Println(diff - 1)
			break
		}
	}
	for diff := 0; ; diff++ {
		loc := lowestCoord
		loc.Y -= diff
		inRange := 0
		for _, d := range drones {
			if d.InRange(loc) {
				inRange++
			}
		}
		if inRange != highScore {
			fmt.Println(diff - 1)
			break
		}
	}
	for diff := 0; ; diff++ {
		loc := lowestCoord
		loc.Z -= diff
		inRange := 0
		for _, d := range drones {
			if d.InRange(loc) {
				inRange++
			}
		}
		if inRange != highScore {
			fmt.Println(diff - 1)
			break
		}
	}
	fmt.Println(lowestCoord.X + lowestCoord.Y + lowestCoord.Z)
}
