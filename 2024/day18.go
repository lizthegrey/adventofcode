package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/lizthegrey/adventofcode/2022/heapq"
)

var inputFile = flag.String("inputFile", "inputs/day18.input", "Relative file path to use as input.")

type coord struct {
	x, y int
}

func (c coord) add(o coord) coord {
	return coord{c.x + o.x, c.y + o.y}
}

const maxCoord = 70

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var incoming []coord
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		incoming = append(incoming, coord{x, y})
	}

	impassable := make(map[coord]bool)
	for i, loc := range incoming {
		impassable[loc] = true
		steps := aStar(impassable, coord{0, 0}, coord{maxCoord, maxCoord})
		if i < 1024 {
			continue
		}
		if i == 1024 {
			fmt.Println(steps)
		}
		if steps == -1 {
			fmt.Printf("%d,%d\n", loc.x, loc.y)
			return
		}
	}
}

func aStar(impassable map[coord]bool, start, target coord) int {
	gScore := map[coord]int{
		start: 0,
	}
	workList := heapq.New[coord]()
	workList.Upsert(start, start.heuristic(target))
	for workList.Len() != 0 {
		// Pop the current node off the worklist.
		current := workList.PopSafe()

		if current == target {
			return gScore[current]
		}
		for _, n := range current.iterate(impassable) {
			proposedScore := gScore[current] + 1
			if previousScore, ok := gScore[n]; !ok || proposedScore < previousScore {
				gScore[n] = proposedScore
				workList.Upsert(n, proposedScore+n.heuristic(target))
			}
		}
	}
	return -1
}

func (c coord) heuristic(target coord) int {
	return (target.x - c.x) + (target.y - c.y)
}

func (s coord) iterate(impassable map[coord]bool) []coord {
	var ret []coord
	for _, step := range []coord{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		next := s.add(step)
		if impassable[next] || next.x < 0 || next.x > maxCoord || next.y < 0 || next.y > maxCoord {
			continue
		}
		ret = append(ret, next)
	}
	return ret
}
