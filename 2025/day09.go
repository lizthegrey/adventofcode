package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day09.input", "Relative file path to use as input.")

type Coord struct {
	X, Y int
}

func (c Coord) Area(o Coord) int {
	return (max(o.X-c.X, c.X-o.X) + 1) * (max(o.Y-c.Y, c.Y-o.Y) + 1)
}

type Outline map[Coord]bool

func (o Outline) paint(from, to Coord) {
	if from.X != to.X && from.Y != to.Y {
		panic("invalid")
	}
	if from.X == to.X {
		for y := min(from.Y, to.Y); y <= max(from.Y, to.Y); y++ {
			o[Coord{from.X, y}] = true
		}
	}
	if from.Y == to.Y {
		for x := min(from.X, to.X); x <= max(from.X, to.X); x++ {
			o[Coord{x, from.Y}] = true
		}
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

	var tiles []Coord
	outline := make(Outline)
	for i, s := range split[:len(split)-1] {
		parts := strings.Split(s, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		cur := Coord{x, y}
		tiles = append(tiles, cur)
		if i > 0 {
			outline.paint(tiles[i-1], tiles[i])
		}
		if i == len(split)-2 {
			outline.paint(tiles[0], tiles[i])
		}
	}

	var highestA, highestB int
	for i, a := range tiles {
	outer:
		for j, b := range tiles {
			if i <= j {
				continue
			}
			area := a.Area(b)
			if area > highestA {
				highestA = area
			}
			if area <= highestB {
				continue
			}
			for k, c := range tiles {
				if k == i || k == j {
					continue
				}
				// Disallow rectangles that contain any other corners in them that create
				// lines that cross our rectangle.
				if min(a.X, b.X) < c.X && c.X < max(a.X, b.X) && min(a.Y, b.Y) < c.Y && c.Y < max(a.Y, b.Y) {
					continue outer
				}
			}
			// Walk along the perimeter of any final candidates, looking for crossings that cut into our area.
			for r, x := range [2]int{min(a.X, b.X), max(a.X, b.X)} {
				if a.Y == b.Y {
					continue
				}
				for y := min(a.Y, b.Y) + 1; y <= max(a.Y, b.Y)-1; y++ {
					if outline[Coord{x + 1 - 2*r, y}] {
						continue outer
					}
				}
			}
			for s, y := range [2]int{min(a.Y, b.Y), max(a.Y, b.Y)} {
				if a.X == b.X {
					continue
				}
				for x := min(a.X, b.X) + 1; x <= max(a.X, b.X)-1; x++ {
					if outline[Coord{x, y + 1 - 2*s}] {
						continue outer
					}
				}
			}
			highestB = area
		}
	}
	fmt.Println(highestA)
	fmt.Println(highestB)
}
