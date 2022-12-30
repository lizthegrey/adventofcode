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

type board struct {
	tiles      map[coord]bool
	maxR, maxC int
}

func (b board) bounds(pos coord, cube bool) coord {
	if pos.r < -1 {
		pos.r = b.maxR + 1
	} else if pos.r > b.maxR+1 {
		pos.r = -1
	}
	if pos.c < -1 {
		pos.c = b.maxC + 1
	} else if pos.c > b.maxC+1 {
		pos.c = -1
	}
	return pos
}

type turtle struct {
	pos    coord
	facing dir
}

func (t *turtle) forward(passable board, n int, cube bool) {
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

	var i int
	rollback := t.pos
	for i < n {
		// Provisionally move us one square with wraparound.
		proposed := t.pos.add(incr)
		proposed = passable.bounds(proposed, cube)

		if p, ok := passable.tiles[proposed]; !ok {
			// This isn't a real tile. slide riiiight on over without incrementing i
			// or changing rollback value
			t.pos = proposed
		} else if !p {
			// We've hit a wall, we can't proceed regardless of having more allowed movement.
			// Roll back to the last known non-void position rather than stranding us in the void.
			t.pos = rollback
			break
		} else {
			i++
			// Change both pos and rollback
			rollback = proposed
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

	var passable board
	passable.tiles = make(map[coord]bool)
	for r, line := range split[:len(split)-3] {
		if r > passable.maxR {
			passable.maxR = r
		}
		for c, v := range line {
			if c > passable.maxC {
				passable.maxC = c
			}
			switch v {
			case ' ':
				continue
			case '#':
				passable.tiles[coord{r, c}] = false
			case '.':
				passable.tiles[coord{r, c}] = true
			}
		}
	}

	// part A
	fmt.Println(passable.run(split[len(split)-2], false))
	// part B
	fmt.Println(passable.run(split[len(split)-2], true))
}

func (b board) run(commands string, cube bool) int {
	player := turtle{}
	player.forward(b, 1, cube)

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
			player.forward(b, v, cube)
			i += j - 1
		}
	}
	return player.password()
}
