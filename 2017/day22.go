package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day22.input", "Relative file path to use as input.")
var cycles = flag.Int("cycles", 5, "The number of cycles to simulate.")

type Coord struct {
	X, Y int
}
type Board map[Coord]bool
type WorkerState struct {
	Environment Board
	Loc         Coord
	Dir         int
}

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes[:len(bytes)-1])
	lines := strings.Split(contents, "\n")

	b := make(Board)

	// always an odd number. coordinates will range from [-height/2..height/2]
	height := len(lines)
	width := len(lines[0])

	// Convert input file into pattern of initial infection.
	for r, l := range lines {
		for c, v := range l {
			if v == '#' {
				loc := Coord{
					X: c - (width / 2),
					Y: (height / 2) - r,
				}
				b[loc] = true
			}
		}
	}

	carrier := WorkerState{b, Coord{0, 0}, 0}

	spread := 0
	for i := 0; i < *cycles; i++ {
		if carrier.Tick() {
			spread++
		}
	}
	fmt.Println(spread)
}

func (ws *WorkerState) Tick() bool {
	infected := false

	// Turn
	if ws.Environment[ws.Loc] {
		ws.Dir = (ws.Dir + 1) % 4
	} else {
		ws.Dir = (ws.Dir + 3) % 4
	}

	// Toggle
	if ws.Environment[ws.Loc] {
		delete(ws.Environment, ws.Loc)
	} else {
		infected = true
		ws.Environment[ws.Loc] = true
	}

	switch ws.Dir {
	case 0:
		ws.Loc = Coord{ws.Loc.X, ws.Loc.Y + 1}
	case 1:
		ws.Loc = Coord{ws.Loc.X + 1, ws.Loc.Y}
	case 2:
		ws.Loc = Coord{ws.Loc.X, ws.Loc.Y - 1}
	case 3:
		ws.Loc = Coord{ws.Loc.X - 1, ws.Loc.Y}
	}

	return infected
}
