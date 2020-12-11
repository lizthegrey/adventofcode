package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day11.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use part B logic.")

type Coord struct {
	Row, Col int
}

func (c Coord) Adjacent(occupied map[Coord]bool) []Coord {
	var ret []Coord
	for rOffset := -1; rOffset <= 1; rOffset++ {
		for cOffset := -1; cOffset <= 1; cOffset++ {
			if cOffset == 0 && rOffset == 0 {
				continue
			}
			if !*partB {
				ret = append(ret, Coord{c.Row + rOffset, c.Col + cOffset})
			} else {
				for i := 1; ; i++ {
					pos := Coord{c.Row + rOffset*i, c.Col + cOffset*i}
					if pos.Col < 0 || pos.Row < 0 || pos.Col > 100 || pos.Row > 100 {
						break
					}
					if _, seat := occupied[pos]; seat {
						ret = append(ret, pos)
						break
					}
				}
			}
		}
	}
	return ret
}

func (c Coord) OccupiedNeighbors(occupied map[Coord]bool) int {
	count := 0
	for _, n := range c.Adjacent(occupied) {
		if occupied[n] {
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
	occupied := make(map[Coord]bool)
	for r, s := range split {
		for c, v := range s {
			switch v {
			case 'L':
				occupied[Coord{r, c}] = false
			case '#':
				occupied[Coord{r, c}] = true
			case '.':
				delete(occupied, Coord{r, c})
			}
		}
	}

	occupyThreshold := 4
	if *partB {
		occupyThreshold = 5
	}
	for {
		changed := false
		prev := make(map[Coord]bool)
		for k, v := range occupied {
			prev[k] = v
		}
		for k, occ := range prev {
			if occ {
				if k.OccupiedNeighbors(prev) >= occupyThreshold {
					occupied[k] = false
					changed = true
				}
			} else {
				if k.OccupiedNeighbors(prev) == 0 {
					occupied[k] = true
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}

	result := 0
	for _, v := range occupied {
		if v {
			result++
		}
	}
	fmt.Println(result)
}
