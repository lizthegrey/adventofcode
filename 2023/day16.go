package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day16.input", "Relative file path to use as input.")

type Coord struct {
	R, C int
}

type Direction int

const (
	North = iota
	East
	South
	West
)

type Beam struct {
	Coord
	Dir Direction
}

func (b Beam) Move() Beam {
	next := b
	switch b.Dir {
	case North:
		next.R--
	case East:
		next.C++
	case South:
		next.R++
	case West:
		next.C--
	}
	return next
}

type Grid map[Coord]rune

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	grid := make(Grid)
	maxC := len(split[0]) - 1
	var maxR int
	for r, s := range split[:len(split)-1] {
		for c, v := range s {
			grid[Coord{r, c}] = v
		}
		maxR = r
	}
	partA := grid.Energize(Beam{Coord{0, 0}, Direction(East)}, maxR, maxC)
	fmt.Println(partA)
	var partB int
	for _, dir := range []Direction{North, East, South, West} {
		var initial Beam
		initial.Dir = dir
		bound := maxC
		if dir == East || dir == West {
			bound = maxR
		}
		for i := 0; i <= bound; i++ {
			switch dir {
			case North:
				initial.R = maxR
				initial.C = i
			case East:
				initial.R = i
				initial.C = 0
			case South:
				initial.R = 0
				initial.C = i
			case West:
				initial.R = i
				initial.C = maxC
			}
			if result := grid.Energize(initial, maxR, maxC); result > partB {
				partB = result
			}
		}
	}
	fmt.Println(partB)
}

func (g Grid) Energize(initial Beam, maxR, maxC int) int {
	active := make(map[Coord]bool)
	seen := make(map[Beam]bool)
	q := []Beam{initial}
	for len(q) > 0 {
		item := q[0]
		q = q[1:]

		if item.R < 0 || item.C < 0 || item.R > maxR || item.C > maxC {
			// Off the board.
			continue
		}
		if seen[item] {
			// Duplicate beam already exists.
			continue
		}
		active[item.Coord] = true
		seen[item] = true
		switch g[item.Coord] {
		case '.':
			q = append(q, item.Move())
		case '/':
			switch item.Dir {
			case North:
				item.Dir = East
			case East:
				item.Dir = North
			case South:
				item.Dir = West
			case West:
				item.Dir = South
			}
			q = append(q, item.Move())
		case '\\':
			switch item.Dir {
			case North:
				item.Dir = West
			case East:
				item.Dir = South
			case South:
				item.Dir = East
			case West:
				item.Dir = North
			}
			q = append(q, item.Move())
		case '|':
			switch item.Dir {
			case North:
				fallthrough
			case South:
				q = append(q, item.Move())
			case East:
				fallthrough
			case West:
				item.Dir = North
				q = append(q, item.Move())
				item.Dir = South
				q = append(q, item.Move())
			}
		case '-':
			switch item.Dir {
			case North:
				fallthrough
			case South:
				item.Dir = East
				q = append(q, item.Move())
				item.Dir = West
				q = append(q, item.Move())
			case East:
				fallthrough
			case West:
				q = append(q, item.Move())
			}
		}
	}
	return len(active)
}
