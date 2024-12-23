package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
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
	shortest, predecessors := dijkstra(passable, start, target)
	fmt.Println(shortest)

	paths := make(map[state]bool)
	for _, dir := range [4]dir{North, East, South, West} {
		backtrace(paths, predecessors, state{target, dir})
	}
	viewing := make(map[coord]bool)
	for k := range paths {
		viewing[k.coord] = true
	}
	fmt.Println(len(viewing))
}

func backtrace(paths map[state]bool, predecessors map[state][]state, cur state) {
	if paths[cur] {
		// Already processed.
		return
	}
	paths[cur] = true
	for _, pre := range predecessors[cur] {
		backtrace(paths, predecessors, pre)
	}
}

func dijkstra(passable map[coord]bool, start, target coord) (int, map[state][]state) {
	initial := state{start, East}
	ret := make(map[state][]state)

	gScore := map[state]int{
		initial: 0,
	}
	workList := heapq.New[state]()
	workList.Upsert(initial, 0)

	targetScore := math.MaxInt
	for workList.Len() != 0 {
		// Pop the current node off the worklist.
		current := workList.PopSafe()
		score := gScore[current]

		// Handle multiple paths to ending condition.
		if score > targetScore {
			return targetScore, ret
		}
		if current.coord == target {
			targetScore = score
			continue
		}

		for _, n := range current.iterate(passable) {
			incr := 1
			if n.coord == current.coord {
				incr = 1000
			}
			proposedScore := score + incr
			if previousScore, ok := gScore[n]; !ok || proposedScore <= previousScore {
				if !ok || proposedScore < previousScore {
					workList.Upsert(n, proposedScore)
					// this is not clear(ret[n]) as that does something different and wrong!
					ret[n] = ret[n][:0]
				}
				ret[n] = append(ret[n], current)
				gScore[n] = proposedScore
			}
		}
	}
	return -1, nil
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
