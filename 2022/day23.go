package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day23.input", "Relative file path to use as input.")

type coord struct {
	r, c int
}

type board map[coord]bool

type rule func(board, coord) *coord

func noop(b board, src coord) bool {
	for r := -1; r <= 1; r++ {
		for c := -1; c <= 1; c++ {
			if r == 0 && c == 0 {
				continue
			}
			if b[coord{src.r + r, src.c + c}] {
				return false
			}
		}
	}
	return true
}

func north(b board, src coord) *coord {
	proposed := coord{src.r - 1, src.c + 0}
	for _, test := range []coord{{proposed.r, src.c - 1}, {proposed.r, src.c + 0}, {proposed.r, src.c + 1}} {
		if b[test] {
			return nil
		}
	}
	return &proposed
}

func south(b board, src coord) *coord {
	proposed := coord{src.r + 1, src.c + 0}
	for _, test := range []coord{{proposed.r, src.c - 1}, {proposed.r, src.c + 0}, {proposed.r, src.c + 1}} {
		if b[test] {
			return nil
		}
	}
	return &proposed
}

func east(b board, src coord) *coord {
	proposed := coord{src.r + 0, src.c + 1}
	for _, test := range []coord{{src.r - 1, proposed.c}, {src.r + 0, proposed.c}, {src.r + 1, proposed.c}} {
		if b[test] {
			return nil
		}
	}
	return &proposed
}

func west(b board, src coord) *coord {
	proposed := coord{src.r + 0, src.c - 1}
	for _, test := range []coord{{src.r - 1, proposed.c}, {src.r + 0, proposed.c}, {src.r + 1, proposed.c}} {
		if b[test] {
			return nil
		}
	}
	return &proposed
}

func (b board) iterate(rules []rule) int {
	// Tracks which elf wants to move to which location.
	proposed := make(map[coord][]coord)

	// Populate desired moves.
	for src := range b {
		if noop(b, src) {
			continue
		}
		// Evaluate the rules in sequence, taking the first match.
		for _, rule := range rules {
			if dst := rule(b, src); dst != nil {
				proposed[*dst] = append(proposed[*dst], src)
				break
			}
		}
	}

	// Actually move the elves, if only one source elf wanted that destination.
	var moved int
	for dst, srcs := range proposed {
		if len(srcs) == 1 {
			delete(b, srcs[0])
			b[dst] = true
			moved += 1
		}
	}
	return moved
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	elves := make(board)
	for r, line := range split[:len(split)-1] {
		for c, v := range line {
			if v == '#' {
				elves[coord{r, c}] = true
			}
		}
	}

	rules := []rule{north, south, west, east}
	for round := 1; ; round++ {
		if elves.iterate(rules) == 0 {
			// Part B
			fmt.Println(round)
			return
		}

		if round == 10 {
			var minR, minC, maxR, maxC int
			for e := range elves {
				if e.r < minR {
					minR = e.r
				}
				if e.r > maxR {
					maxR = e.r
				}
				if e.c < minC {
					minC = e.c
				}
				if e.c > maxC {
					maxC = e.c
				}
			}
			height := maxR - minR + 1
			width := maxC - minC + 1
			occupied := len(elves)
			// Part A
			fmt.Println(height*width - occupied)
		}

		head := rules[0]
		rules = append(rules[1:], head)
	}
}
