package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/lizthegrey/adventofcode/2022/heapq"
)

var inputFile = flag.String("inputFile", "inputs/day16.input", "Relative file path to use as input.")

type coord struct {
	r, c int
}

type dir uint8

const (
	North dir = iota
	East
	South
	West
)

type state struct {
	coord
	facing dir
}

type augmented struct {
	state
	steps int
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var start, target coord
	passable := make(map[coord]bool)
	for r, s := range split[:len(split)-1] {
		for c, v := range s {
			loc := coord{r, c}
			switch v {
			case '.':
				passable[loc] = true
			case '#':
			case 'S':
				passable[loc] = true
				start = loc
			case 'E':
				passable[loc] = true
				target = loc
			}
		}
	}
	shortest := aStar(passable, start, target)
	fmt.Println(shortest)
	turns := shortest / 1000
	steps := shortest % 1000

	watch := make(map[coord]bool)
	// Having constrained the maze, now we can enumerate all paths from start
	// that reach finish in exactly the number of turns and steps, and not one more or less.
	path := []state{{start, East}}
	visited := make(map[coord]bool)
	memo := make(map[state]augmented)
	explore(passable, visited, watch, memo, path, target, turns, steps)
	fmt.Println(len(watch))
}

func explore(passable, visited, watch map[coord]bool, memo map[state]augmented, path []state, target coord, turns, steps int) {
	current := path[len(path)-1]
	if current.coord == target {
		for i, loc := range path {
			watch[loc.coord] = true
			if i > 0 {
				prev := path[i-1].coord
				next := loc.coord
				for r := min(prev.r, next.r); r <= max(prev.r, next.r); r++ {
					watch[coord{r, prev.c}] = true
				}
				for c := min(prev.c, next.c); c <= max(prev.c, next.c); c++ {
					watch[coord{prev.r, c}] = true
				}
			}
		}
		fmt.Println(len(watch))
		return
	}
	visited[current.coord] = true
	for _, n := range current.iterate(passable) {
		turnsLeft := turns
		stepsLeft := steps
		if n.coord == current.coord {
			turnsLeft--
			if turnsLeft < 0 {
				continue
			}
		} else {
			if zoom, ok := memo[n]; ok {
				n = zoom.state
				stepsLeft -= zoom.steps
			} else {
				orig := n
				i := 1
				for {
					options := n.iterate(passable)
					if len(options) != 1 || options[0].facing != n.facing {
						break
					}
					n = options[0]
					i++
				}
				memo[orig] = augmented{n, i}
				stepsLeft -= i
			}
			if stepsLeft < 0 {
				// todo: also account for stepping directly to end, if we don't have enough steps left to do that it's futile.
				continue
			}
			if visited[n.coord] {
				continue
			}
		}
		path = append(path, n)
		explore(passable, visited, watch, memo, path, target, turnsLeft, stepsLeft)
		path = path[:len(path)-1]
	}
	delete(visited, current.coord)
}

func aStar(passable map[coord]bool, start, target coord) int {
	initial := state{start, East}

	gScore := map[state]int{
		initial: 0,
	}
	workList := heapq.New[state]()
	workList.Upsert(initial, start.heuristic(target))
	for workList.Len() != 0 {
		// Pop the current node off the worklist.
		current := workList.PopSafe()

		if current.coord == target {
			return gScore[current]
		}
		for _, n := range current.iterate(passable) {
			var incr int
			if n.coord == current.coord {
				incr += 1000
			} else {
				incr++
			}
			proposedScore := gScore[current] + incr
			if previousScore, ok := gScore[n]; !ok || proposedScore < previousScore {
				gScore[n] = proposedScore
				workList.Upsert(n, proposedScore+n.heuristic(target))
			}
		}
	}
	return -1
}

func (c coord) heuristic(target coord) int {
	return (target.r - c.r) + (target.c - c.c)
}

func (s state) iterate(passable map[coord]bool) []state {
	var ret []state
	straight := s
	switch s.facing {
	case North:
		straight.r--
	case East:
		straight.c++
	case South:
		straight.r++
	case West:
		straight.c--
	}
	if passable[straight.coord] {
		ret = append(ret, straight)
	}
	for _, dir := range [4]dir{North, East, South, West} {
		if dir%2 == s.facing%2 {
			// Can only turn in 90 degree increments
			continue
		}
		cpy := s
		cpy.facing = dir
		next := s.coord
		switch cpy.facing {
		case North:
			next.r--
		case East:
			next.c++
		case South:
			next.r++
		case West:
			next.c--
		}
		if !passable[next] {
			continue
		}
		ret = append(ret, cpy)
	}
	return ret
}
