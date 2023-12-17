package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/lizthegrey/adventofcode/2022/heapq"
)

var inputFile = flag.String("inputFile", "inputs/day17.input", "Relative file path to use as input.")

type Coord struct {
	R, C int
}

type Direction int

const (
	Stationary = iota
	Up
	Right
	Down
	Left
)

type Step struct {
	Coord
	Dir      Direction
	Momentum int
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	grid := make(map[Coord]int)
	maxC := len(split[0]) - 1
	var maxR int
	for r, row := range split[:len(split)-1] {
		for c, v := range row {
			grid[Coord{r, c}] = int(v - '0')
		}
		maxR = r
	}
	start := Coord{0, 0}
	target := Coord{maxR, maxC}
	fmt.Println(aStar(grid, start, target, false))
	fmt.Println(aStar(grid, start, target, true))
}

func aStar(g map[Coord]int, source, target Coord, ultra bool) int {
	initial := Step{source, Stationary, 0}
	gScore := map[Step]int{
		initial: 0,
	}
	workList := heapq.New[Step]()
	workList.Upsert(initial, source.Heuristic(target))
	for workList.Len() != 0 {
		// Pop the current node off the worklist.
		current := workList.PopSafe()

		if current.Coord == target && (!ultra || current.Momentum >= 4) {
			return gScore[current]
		}
		for _, n := range current.Iterate(target, ultra) {
			proposedScore := gScore[current] + g[n.Coord]
			if previousScore, ok := gScore[n]; !ok || proposedScore < previousScore {
				gScore[n] = proposedScore
				workList.Upsert(n, proposedScore+n.Heuristic(target))
			}
		}
	}
	return -1
}

func (c Coord) Heuristic(target Coord) int {
	return (target.R - c.R) + (target.C - c.C)
}

func (s Step) Iterate(target Coord, ultra bool) []Step {
	var ret []Step
	for _, dir := range [4]Direction{Up, Right, Down, Left} {
		cpy := s
		cpy.Dir = dir
		if dir != s.Dir {
			if ultra && s.Dir != Stationary && s.Momentum < 4 {
				continue
			}
			cpy.Momentum = 0
		}
		cpy.Momentum++
		if !ultra && cpy.Momentum > 3 {
			continue
		}
		if ultra && cpy.Momentum > 10 {
			continue
		}

		switch dir {
		case Up:
			if s.Dir == Down {
				continue
			}
			cpy.R--
		case Right:
			if s.Dir == Left {
				continue
			}
			cpy.C++
		case Down:
			if s.Dir == Up {
				continue
			}
			cpy.R++
		case Left:
			if s.Dir == Right {
				continue
			}
			cpy.C--
		}
		if cpy.R < 0 || cpy.R > target.R || cpy.C < 0 || cpy.C > target.C {
			continue
		}
		ret = append(ret, cpy)
	}
	return ret
}
