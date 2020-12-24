package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day24.input", "Relative file path to use as input.")

// Coord is a 2-axis representation of the hex grid.
//    (-1, 1)  (0, 1)  (1, 1)
// (-1,0)  (0,0)  (1, 0)
//    (0,-1)  (1, -1)
type Coord struct {
	X, Y int
}

func (c Coord) Neighbors() []Coord {
	return []Coord{
		{c.X + 1, c.Y},
		{c.X - 1, c.Y},
		{c.X, c.Y + 1},
		{c.X, c.Y - 1},
		{c.X + 1, c.Y - 1},
		{c.X - 1, c.Y + 1},
	}
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

	black := make(map[Coord]bool)
	for _, s := range split {
		current := Coord{0, 0}
		for i := 0; i < len(s); i++ {
			switch s[i] {
			case 'e':
				current.X++
			case 'w':
				current.X--
			case 'n':
				switch s[i+1] {
				case 'e':
					current.Y++
				case 'w':
					current.X--
					current.Y++
				default:
					fmt.Printf("Failed to parse value %s\n", s)
				}
				i++
			case 's':
				switch s[i+1] {
				case 'e':
					current.X++
					current.Y--
				case 'w':
					current.Y--
				default:
					fmt.Printf("Failed to parse value %s\n", s)
				}
				i++
			default:
				fmt.Printf("Failed to parse value %s\n", s)
			}
		}
		if black[current] {
			delete(black, current)
		} else {
			black[current] = true
		}
	}
	fmt.Println(len(black))

	for d := 0; d < 100; d++ {
		workingSet := make(map[Coord]bool)
		for t := range black {
			workingSet[t] = true
			for _, n := range t.Neighbors() {
				if _, found := workingSet[n]; !found {
					workingSet[n] = false
				}
			}
		}
		next := make(map[Coord]bool)
		for t, b := range workingSet {
			neighbors := t.Neighbors()
			neighborCount := 0
			for _, n := range neighbors {
				if black[n] {
					neighborCount++
				}
			}
			if b {
				if neighborCount == 1 || neighborCount == 2 {
					next[t] = true
				}
			} else {
				if neighborCount == 2 {
					next[t] = true
				}
			}
		}
		black = next
	}
	fmt.Println(len(black))
}
