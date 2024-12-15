package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day15.input", "Relative file path to use as input.")

type coord struct {
	r, c int
}

func (c coord) add(o coord) coord {
	return coord{c.r + o.r, c.c + o.c}
}

const (
	Up dir = iota
	Down
	Left
	Right
)

type dir uint8

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	wallsA := make(map[coord]bool)
	boxesA := make(map[coord]bool)
	wallsB := make(map[coord]bool)
	boxesB := make(map[coord]bool)
	var robotA, robotB coord
	var mazeFinished bool
	var instrs []dir
	for r, s := range split[:len(split)-1] {
		if len(s) == 0 {
			mazeFinished = true
			continue
		}
		if mazeFinished {
			for _, v := range s {
				var instr dir
				switch v {
				case '<':
					instr = Left
				case '>':
					instr = Right
				case '^':
					instr = Up
				case 'v':
					instr = Down
				}
				instrs = append(instrs, instr)
			}
			continue
		}
		for c, v := range s {
			loc := coord{r, c}
			switch v {
			case '#':
				wallsA[loc] = true
				wallsB[coord{r, 2 * c}] = true
				wallsB[coord{r, 2*c + 1}] = true
			case '.':
			case 'O':
				boxesA[loc] = true
				boxesB[coord{r, 2 * c}] = true
			case '@':
				robotA = loc
				robotB = coord{r, 2 * c}
			default:
				panic("Found unknown character in input")
			}
		}
	}

	// part A
	for _, instr := range instrs {
		var increment coord
		switch instr {
		case Left:
			increment.c--
		case Right:
			increment.c++
		case Up:
			increment.r--
		case Down:
			increment.r++
		}
		proposed := robotA.add(increment)
		if (!wallsA[proposed] && !boxesA[proposed]) || shuffleA(wallsA, boxesA, proposed, increment) {
			robotA = proposed
		}
	}

	var scoreA int
	for loc := range boxesA {
		scoreA += 100*loc.r + loc.c
	}
	fmt.Println(scoreA)

	// part B
	for _, instr := range instrs {
		var increment coord
		switch instr {
		case Left:
			increment.c--
		case Right:
			increment.c++
		case Up:
			increment.r--
		case Down:
			increment.r++
		}
		proposed := robotB.add(increment)
		box := normaliseBox(boxesB, proposed)
		if !wallsB[proposed] && box == nil {
			robotB = proposed
		} else if shuffleB(wallsB, boxesB, proposed, increment, true) {
			shuffleB(wallsB, boxesB, proposed, increment, false)
			robotB = proposed
		}
	}

	var scoreB int
	for loc := range boxesB {
		scoreB += 100*loc.r + loc.c
	}
	fmt.Println(scoreB)
}

// shuffle returns true if it's possible to move, false (and no change) if it is not possible to move.
func shuffleA(walls, boxes map[coord]bool, proposed, increment coord) bool {
	if walls[proposed] {
		// Walls cannot be moved.
		return false
	}
	if !boxes[proposed] {
		// This is an empty space that we can move a box or robot into.
		return true
	}
	// At this point we know there is a box in the way. see if we can chain move.
	target := proposed.add(increment)
	if shuffleA(walls, boxes, target, increment) {
		delete(boxes, proposed)
		boxes[target] = true
		return true
	}
	return false
}

func normaliseBox(boxes map[coord]bool, loc coord) *coord {
	ret := loc
	if boxes[ret] {
		return &ret
	}
	ret.c--
	if boxes[ret] {
		return &ret
	}
	return nil
}

// shuffleB returns true if it's possible to move, false (and no change) if it is not possible to move.
// has a dry run param to test rather than carry out moves.
func shuffleB(walls, boxes map[coord]bool, proposed, increment coord, dryRun bool) bool {
	if walls[proposed] {
		// Walls cannot be moved.
		return false
	}
	box := normaliseBox(boxes, proposed)
	if box == nil {
		// This is an empty space that we can move a box or robot into.
		return true
	}
	left := *box
	right := coord{left.r, left.c + 1}
	// At this point we know there is a box in the way. see if we can chain move.
	// Left or right moves do not fork, but do require us to consider the correct square.
	// Up/down moves do potentially fork.
	ok := true
	switch increment.c {
	case -1:
		ok = shuffleB(walls, boxes, left.add(increment), increment, dryRun)
	case 1:
		ok = shuffleB(walls, boxes, right.add(increment), increment, dryRun)
	case 0:
		ok = shuffleB(walls, boxes, left.add(increment), increment, dryRun) && shuffleB(walls, boxes, right.add(increment), increment, dryRun)
	}
	if ok {
		if !dryRun {
			delete(boxes, left)
			boxes[left.add(increment)] = true
		}
		return true
	}
	if !dryRun {
		panic("instructed to do a move that failed")
	}
	return false
}
