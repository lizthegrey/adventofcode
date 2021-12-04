package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day04.input", "Relative file path to use as input.")

type Board struct {
	Numbers [5][5]int
	Marked  [5][5]bool
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

	numbers := strings.Split(split[0], ",")
	called := make([]int, 0)
	for _, n := range numbers {
		i, _ := strconv.Atoi(n)
		called = append(called, i)
	}
	boards := make([]*Board, 0)
	for i := 2; i < len(split); i += 6 {
		boards = append(boards, parseBoard(split[i:i+5]))
	}

outer:
	for _, call := range called {
		for _, b := range boards {
			b.mark(call)
			if b.bingo() {
				fmt.Println(b.score(call))
				break outer
			}
		}
	}

	boards = make([]*Board, 0)
	for i := 2; i < len(split); i += 6 {
		boards = append(boards, parseBoard(split[i:i+5]))
	}
	won := make(map[int]bool)
partB:
	for _, call := range called {
		for i, b := range boards {
			if won[i] {
				continue
			}
			b.mark(call)
			if b.bingo() {
				won[i] = true
				if len(won) == len(boards) {
					fmt.Println(b.score(call))
					break partB
				}
			}
		}
	}
}

func parseBoard(lines []string) *Board {
	var b Board
	for r, l := range lines {
		c := 0
		for pos := 0; pos < len(l); pos += 3 {
			n, _ := strconv.Atoi(strings.Trim(l[pos:pos+2], " "))
			b.Numbers[r][c] = n
			c++
		}
	}
	return &b
}

func (b *Board) mark(v int) {
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if b.Numbers[r][c] == v {
				b.Marked[r][c] = true
			}
		}
	}
}

func (b *Board) bingo() bool {
	// Check horizontal bingos
rows:
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if !b.Marked[r][c] {
				continue rows
			}
		}
		return true
	}
	// Check vertical bingos
columns:
	for c := 0; c < 5; c++ {
		for r := 0; r < 5; r++ {
			if !b.Marked[r][c] {
				continue columns
			}
		}
		return true
	}
	return false
}

func (b Board) score(final int) int {
	sum := 0
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if !b.Marked[r][c] {
				sum += b.Numbers[r][c]
			}
		}
	}
	return sum * final
}
