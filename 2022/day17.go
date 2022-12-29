package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day17.input", "Relative file path to use as input.")

const width = 7

type coord struct {
	x, y int
}

type longs struct {
	x, y int64
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

type terrain map[longs]bool

func (t terrain) collides(s stencil, xOff, yOff int64) bool {
	if yOff < 0 {
		return true
	}
	for _, c := range s {
		x := int64(c.x) + xOff
		if x < 0 || x >= width || t[longs{x, int64(c.y) + yOff}] {
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

func (s sequence[T]) modulus() int {
	return s.n
}

func (s sequence[T]) length() int {
	return len(s.elems)
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
	fmt.Println(simulate(shapes, moves, 1000000000000))
}

func simulate(shapes sequence[stencil], moves sequence[bool], rounds int64) int64 {
	var maxY int64
	board := make(terrain)
	cache := make(map[coord]longs)
	var streak int
	for i := int64(0); i < rounds; i++ {
		if entry, ok := cache[coord{shapes.modulus(), moves.modulus()}]; ok {
			streak += 1
			if streak == shapes.length() {
				// We've started over. Fast forward it.
				cycleLen := i - entry.x
				heightDelta := maxY - entry.y
				repeats := (rounds - i) / cycleLen
				i += cycleLen * repeats

				yOffset := heightDelta * repeats
				for y := entry.y; y <= maxY; y++ {
					for x := 0; x < width; x++ {
						if board[longs{int64(x), y}] {
							board[longs{int64(x), y + yOffset}] = true
						}
					}
				}
				maxY += yOffset
			}
		} else {
			// Streak has been reset.
			streak = 0
			cache[coord{shapes.modulus(), moves.modulus()}] = longs{i, maxY}
		}

		shape := shapes.next()
		xOff := int64(2)
		yOff := maxY + int64(3)
		for {
			// Attempt to move the piece sideways,
			left := moves.next()
			if left {
				if !board.collides(shape, xOff-int64(1), yOff) {
					xOff -= 1
				}
			} else {
				if !board.collides(shape, xOff+int64(1), yOff) {
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
			board[longs{int64(c.x) + xOff, int64(c.y) + yOff}] = true
		}
		top := yOff + int64(shape.height())
		if top > maxY {
			maxY = top
		}
	}
	return maxY
}
