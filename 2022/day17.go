package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day17.input", "Relative file path to use as input.")

type coord struct {
	x, y int
}

type stencil []coord

func (s stencil) height() int {
	var max int
	for _, c := range s {
		if c.y > max {
			max = c.y
		}
	}
	return 1 + max
}

type terrain map[coord]bool

func (t terrain) collides(s stencil, xOff, yOff int) bool {
	if yOff < 0 {
		return true
	}
	for _, c := range s {
		if c.x+xOff < 0 || c.x+xOff >= 7 || t[coord{c.x + xOff, c.y + yOff}] {
			return true
		}
	}
	return false
}

type sequence[T any] struct {
	elems []T
	n     int
}

func (s *sequence[T]) next() T {
	n := s.n
	s.n = (s.n + 1) % len(s.elems)
	return s.elems[n]
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var moves sequence[bool]
	for _, s := range split[0] {
		switch s {
		case '>':
			moves.elems = append(moves.elems, false)
		case '<':
			moves.elems = append(moves.elems, true)
		}
	}

	shapes := sequence[stencil]{elems: []stencil{
		{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
		{{1, 0}, {0, 1}, {1, 1}, {2, 1}, {1, 2}},
		{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}},
		{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
		{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
	}}

	// part A
	fmt.Println(simulate(shapes, moves, 2022))
	// part B
	// We want to know after 1000000000000 pieces.
	// This isn't practical to compute by hand, instead we need to look for the repeating pattern.
}

func simulate(shapes sequence[stencil], moves sequence[bool], rounds int) int {
	var maxY int
	board := make(terrain)
	for i := 0; i < rounds; i++ {
		shape := shapes.next()
		xOff := 2
		yOff := maxY + 3
		for {
			// Attempt to move the piece sideways,
			left := moves.next()
			if left {
				if !board.collides(shape, xOff-1, yOff) {
					xOff -= 1
				}
			} else {
				if !board.collides(shape, xOff+1, yOff) {
					xOff += 1
				}
			}
			// then attempt to move the piece down.
			if board.collides(shape, xOff, yOff-1) {
				break
			}
			yOff -= 1
		}
		// Land the piece.
		for _, c := range shape {
			board[coord{c.x + xOff, c.y + yOff}] = true
		}
		top := yOff + shape.height()
		if top > maxY {
			maxY = top
		}
	}
	return maxY
}
