package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day17.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use part B logic.")

type Coord4 struct {
	W, X, Y, Z int
}

func (c Coord4) Adjacent() []Coord4 {
	var ret []Coord4
	for xOffset := -1; xOffset <= 1; xOffset++ {
		for yOffset := -1; yOffset <= 1; yOffset++ {
			for zOffset := -1; zOffset <= 1; zOffset++ {
				for wOffset := -1; wOffset <= 1; wOffset++ {
					if wOffset != 0 && !*partB {
						continue
					}
					if xOffset == 0 && yOffset == 0 && zOffset == 0 && wOffset == 0 {
						continue
					}
					ret = append(ret, Coord4{c.W + wOffset, c.X + xOffset, c.Y + yOffset, c.Z + zOffset})
				}
			}
		}
	}
	return ret
}

func (c Coord4) ActiveNeighbors(active map[Coord4]bool) int {
	count := 0
	for _, n := range c.Adjacent() {
		if active[n] {
			count++
		}
	}
	return count
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]
	active := make(map[Coord4]bool)
	for y, s := range split {
		for x, v := range s {
			switch v {
			case '#':
				active[Coord4{0, x, y, 0}] = true
			}
		}
	}

	for i := 0; i < 6; i++ {
		prev := make(map[Coord4]bool)
		for k := range active {
			for _, n := range k.Adjacent() {
				if _, exists := prev[n]; !exists {
					prev[n] = false
				}
			}
			prev[k] = true
		}

		for k, act := range prev {
			if act {
				activeNeighbors := k.ActiveNeighbors(prev)
				if activeNeighbors != 2 && activeNeighbors != 3 {
					delete(active, k)
				}
			} else {
				if k.ActiveNeighbors(prev) == 3 {
					active[k] = true
				}
			}
		}
	}

	result := 0
	for _, v := range active {
		if v {
			result++
		}
	}
	fmt.Println(result)
}
