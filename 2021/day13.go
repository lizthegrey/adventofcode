package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day13.input", "Relative file path to use as input.")

type Coord struct {
	X, Y int
}
type Paper map[Coord]bool
type Fold struct {
	XAxis  bool
	Offset int
}

func (f Fold) Perform(p Paper) Paper {
	newPaper := make(Paper)
	for coord := range p {
		var diff int
		if f.XAxis {
			diff = coord.X - f.Offset
		} else {
			diff = coord.Y - f.Offset
		}
		if diff > 0 {
			// Mirror it over the Offset line.
			if f.XAxis {
				coord.X = f.Offset - diff
			} else {
				coord.Y = f.Offset - diff
			}
			newPaper[coord] = true
		} else if diff < 0 {
			// Retain its current position.
			newPaper[coord] = true
		} else {
			// We should never fold on a dot.
			fmt.Printf("Attempted to fold on X=%d, but dot on line.\n", f.Offset)
		}
	}
	return newPaper
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

	parsingCoordinates := true
	paper := make(Paper)
	var folds []Fold
	for _, s := range split {
		if s == "" {
			parsingCoordinates = false
			continue
		}
		if parsingCoordinates {
			parts := strings.Split(s, ",")
			x, _ := strconv.Atoi(parts[0])
			y, _ := strconv.Atoi(parts[1])
			paper[Coord{x, y}] = true
		} else {
			// fold along x=655
			var fold Fold
			if s[len("fold along ")] == byte('x') {
				fold.XAxis = true
			}
			fold.Offset, _ = strconv.Atoi(s[len("fold along x="):])
			folds = append(folds, fold)
		}
	}
	for i, f := range folds {
		paper = f.Perform(paper)
		if i == 0 {
			fmt.Println(len(paper))
		}
	}
	var maxX, maxY int
	for coords := range paper {
		if coords.X > maxX {
			maxX = coords.X
		}
		if coords.Y > maxY {
			maxY = coords.Y
		}
	}
	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			if paper[Coord{x, y}] {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
}
