package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day03.input", "Relative file path to use as input.")

type Coord struct {
	X, Y int
}

type Tracker [2]float64

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	hit := make(map[Coord]Tracker)
	for i, s := range split {
		if s == "" {
			continue
		}
		loc := Coord{0, 0}
		traversed := float64(0)
		for _, step := range strings.Split(s, ",") {
			distance, err := strconv.Atoi(step[1:])
			if err != nil {
				fmt.Printf("Failed to parse step %s\n", step)
				return
			}
			var f func(*Coord)
			switch step[0] {
			case 'R':
				f = func(c *Coord) {
					c.X++
				}
			case 'L':
				f = func(c *Coord) {
					c.X--
				}
			case 'U':
				f = func(c *Coord) {
					c.Y++
				}
			case 'D':
				f = func(c *Coord) {
					c.Y--
				}
			default:
				fmt.Printf("Failed to parse step %s\n", step)
			}
			for n := 1; n <= distance; n++ {
				f(&loc)
				hits := hit[loc]
				if hits[i] != float64(0) {
					hits[i] = math.Min(hits[i], float64(n)+traversed)
				} else {
					hits[i] = float64(n) + traversed
				}
				hit[loc] = hits
			}
			traversed += float64(distance)
		}
	}
	closest := math.MaxFloat64
	least := math.MaxFloat64
	for loc, hits := range hit {
		if hits[0] != float64(0) && hits[1] != float64(0) {
			distance := math.Abs(float64(loc.X)) + math.Abs(float64(loc.Y))
			if distance < closest {
				closest = distance
			}
			sum := hits[0] + hits[1]
			if sum < least {
				least = sum
			}
		}
	}
	fmt.Printf("Part A: %d, Part B: %d\n", int(closest), int(least))
}
