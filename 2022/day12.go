package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/lizthegrey/adventofcode/2022/heapq"
)

var inputFile = flag.String("inputFile", "inputs/day12.input", "Relative file path to use as input.")

type coord struct {
	x, y int
}

type terrain map[coord]int

func (c coord) minDst(o coord) int {
	var sum int
	if c.x > o.x {
		sum += c.x - o.x
	} else {
		sum += o.x - c.x
	}
	if c.y > o.y {
		sum += c.y - o.y
	} else {
		sum += o.y - c.y
	}
	return sum
}

func aStar(r terrain, src, dst coord) int {
	gScore := terrain{
		src: 0,
	}
	workList := heapq.New[coord]()
	workList.Upsert(src, src.minDst(dst))

	for workList.Len() != 0 {
		// Pop the current node off the worklist.
		current := workList.PopSafe()

		if current == dst {
			return gScore[dst]
		}
		for _, n := range r.neighbours(current) {
			proposedScore := gScore[current] + 1
			if previousScore, ok := gScore[n]; !ok || proposedScore < previousScore {
				gScore[n] = proposedScore
				workList.Upsert(n, proposedScore + n.minDst(dst))
			}
		}
	}
	return -1
}

func (t terrain) canMove(src, dst coord) bool {
	if _, ok := t[dst]; !ok {
		// Don't allow leaving the board.
		return false
	}
	// Allow moving up at most one, but down as much as you like.
	return t[dst]-t[src] <= 1
}

func (t terrain) maybeAppend(list []coord, src, dst coord) []coord {
	if t.canMove(src, dst) {
		list = append(list, dst)
	}
	return list
}

func (t terrain) neighbours(src coord) []coord {
	var ret []coord
	ret = t.maybeAppend(ret, src, coord{src.x + 0, src.y + 1})
	ret = t.maybeAppend(ret, src, coord{src.x + 0, src.y - 1})
	ret = t.maybeAppend(ret, src, coord{src.x + 1, src.y + 0})
	ret = t.maybeAppend(ret, src, coord{src.x - 1, src.y + 0})
	return ret
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var src, dst coord
	env := make(terrain)
	for y, s := range split[:len(split)-1] {
		for x, c := range s {
			loc := coord{x, y}
			var height int
			switch c {
			case 'S':
				height = 0
				src = loc
			case 'E':
				height = 25
				dst = loc
			default:
				height = int(c - 'a')
			}
			env[loc] = height
		}
	}

	// part A
	// Perform an A* search with a worklist.
	steps := aStar(env, src, dst)
	fmt.Println(steps)

	// part B
	// Perform a repeated search.
	minSteps := steps
	for k, v := range env {
		if v == 0 {
			attempt := aStar(env, k, dst)
			if attempt >= 0 && attempt < minSteps {
				minSteps = attempt
			}
		}
	}
	fmt.Println(minSteps)
}
