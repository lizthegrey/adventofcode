package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day06.input", "Relative file path to use as input.")

type coord struct {
	r, c int
}

type state struct {
	coord
	facing int
}

const (
	Up int = iota
	Right
	Down
	Left
)

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var start coord
	passable := make(map[coord]bool)
	for r, s := range split[:len(split)-1] {
		for c, v := range s {
			loc := coord{r, c}
			passable[loc] = v != '#'
			if v == '^' {
				start = loc
			}
		}
	}
	seen := check(passable, start, nil)
	fmt.Println(len(seen))

	var valid int
	for obs := range seen {
		if obs == start {
			continue
		}
		if check(passable, start, &obs) == nil {
			valid++
		}
	}
	fmt.Println(valid)
}

func check(passable map[coord]bool, start coord, obs *coord) map[coord]bool {
	cur := start
	facing := Up
	seen := make(map[coord]bool)
	states := make(map[state]bool)
	for {
		seen[cur] = true
		key := state{cur, facing}
		if states[key] {
			return nil
		}
		states[key] = true
		next := cur
		switch facing {
		case Up:
			next.r--
		case Right:
			next.c++
		case Down:
			next.r++
		case Left:
			next.c--
		}
		if p, ok := passable[next]; !ok {
			// Off board.
			break
		} else if !p || (obs != nil && *obs == next) {
			facing++
			facing %= 4
		} else {
			cur = next
		}
	}
	return seen
}
