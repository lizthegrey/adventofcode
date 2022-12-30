package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day22.input", "Relative file path to use as input.")

type dir int

const (
	right dir = iota
	down
	left
	up
)

type coord struct {
	r, c int
}

// Generates a new result.
func (n coord) add(o coord) coord {
	n.r += o.r
	n.c += o.c
	return n
}

type board map[coord]bool

type turtle struct {
	pos    coord
	facing dir
}

func (t *turtle) forward(passable board, n int) {
	var incr coord
	switch t.facing {
	case right:
		incr.c += 1
	case down:
		incr.r += 1
	case left:
		incr.c -= 1
	case up:
		incr.r -= 1
	}

	// Insert something here for handling n
	var i int
	loc = t.pos
	for {
		// Provisionally move us one square.
		proposed := loc.add(incr)
		// Handle the wall collision and wrapping logic here.
		if p, ok := passable[proposed]; !ok {
			// This isn't a real tile. slide riiiight on over without incrementing i
		} else if i == n || !p {
			// We've hit a wall. Regardless of what n says, we can't proceed any further.
			t.pos = proposed
			break
		} else {
			t.pos = proposed
		}
	}
}

func (t turtle) password() int {
	sum := 1000 * (t.pos.r + 1)
	sum += 4 * (t.pos.c + 1)
	sum += int(t.facing)
	return sum
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	passable := make(board)
	for r, line := range split[:len(split)-3] {
		for c, v := range line {
			switch v {
			case ' ':
				continue
			case '#':
				passable[coord{r, c}] = false
			case '.':
				passable[coord{r, c}] = true
			}
		}
	}
	fmt.Println(len(passable))

	// part A
	// Evaluate the command string now that we have our board.
	player := turtle{}
	player.forward(passable, 1)

	commands := split[len(split)-2]
	for i := 0; i < len(commands); i++ {
		switch commands[i] {
		case 'L':
			player.facing = (player.facing + 3) % 4
		case 'R':
			player.facing = (player.facing + 1) % 4
		default:
			// This is a numeric value that we need to fast forward until we encounter a non-numeric value.
			j := 1
			for i+j < len(commands) && commands[i+j] != 'L' && commands[i+j] != 'R' {
				j++
			}
			v, _ := strconv.Atoi(commands[i : i+j])
			player.forward(passable, v)
			i += j - 1
		}
	}
	fmt.Println(player.password())
}
