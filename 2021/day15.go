package main

import (
	"container/heap"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/lizthegrey/adventofcode/2021/trace"
	"go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/attribute"
)

var inputFile = flag.String("inputFile", "inputs/day15.input", "Relative file path to use as input.")

var tr = otel.Tracer("day15")

type RiskMap map[Coord]int
type Coord struct {
	R, C int
}
type HeapQueue struct {
	Elems     *[]Coord
	Score     RiskMap
	Positions map[Coord]int
}

func (h HeapQueue) Len() int           { return len(*h.Elems) }
func (h HeapQueue) Less(i, j int) bool { return h.Score[(*h.Elems)[i]] < h.Score[(*h.Elems)[j]] }
func (h HeapQueue) Swap(i, j int) {
	h.Positions[(*h.Elems)[i]], h.Positions[(*h.Elems)[j]] = h.Positions[(*h.Elems)[j]], h.Positions[(*h.Elems)[i]]
	(*h.Elems)[i], (*h.Elems)[j] = (*h.Elems)[j], (*h.Elems)[i]
}

func (h HeapQueue) Push(x interface{}) {
	h.Positions[x.(Coord)] = len(*h.Elems)
	*h.Elems = append(*h.Elems, x.(Coord))
}

func (h HeapQueue) Pop() interface{} {
	old := *h.Elems
	n := len(old)
	x := old[n-1]
	*h.Elems = old[0 : n-1]
	delete(h.Positions, x)
	return x
}

func (h HeapQueue) Position(x Coord) int {
	if pos, ok := h.Positions[x]; ok {
		return pos
	}
	return -1
}

func main() {
	flag.Parse()

	ctx := context.Background()
	hny, tp := trace.InitializeTracing(ctx)
	defer hny.Shutdown(ctx)
	defer tp.Shutdown(ctx)

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	risks := make(RiskMap)
	for r, line := range split {
		for c, risk := range line {
			risks[Coord{r, c}] = int(risk - '0')
		}
	}
	height := len(split)
	width := len(split[0])
	start := Coord{0, 0}
	dst := Coord{height - 1, width - 1}
	fmt.Println(AStar(risks, &start, &dst))

	dst = Coord{height*5 - 1, width*5 - 1}
	expandedRisks := make(RiskMap)
	for k, v := range risks {
		for r := 0; r < 5; r++ {
			for c := 0; c < 5; c++ {
				increase := r + c
				value := 1 + (v+increase-1)%9
				expandedRisks[Coord{k.R + r*height, k.C + c*width}] = value
			}
		}
	}
	fmt.Println(AStar(expandedRisks, &start, &dst))
}

func AStar(r RiskMap, src, dst *Coord) int {
	gScore := RiskMap{
		*src: 0,
	}
	fScore := RiskMap{
		*src: src.Heuristic(dst),
	}
	workList := HeapQueue{&[]Coord{*src}, fScore, make(RiskMap)}
	heap.Init(&workList)
	history := make(map[Coord]Coord)

	for len(*workList.Elems) != 0 {
		// Pop the current node off the worklist.
		current := heap.Pop(&workList).(Coord)

		if current == *dst {
			// Reconstruct the score by retracing our path to start.
			score := 0
			for current != *src {
				score += r[current]
				current = history[current]
			}
			return score
		}
		for _, n := range r.Neighbors(current) {
			proposedScore := gScore[current] + r[n]
			if previousScore, ok := gScore[n]; !ok || proposedScore < previousScore {
				history[n] = current
				gScore[n] = proposedScore
				fScore[n] = proposedScore + n.Heuristic(dst)
				if pos := workList.Position(n); pos == -1 {
					heap.Push(&workList, n)
				} else {
					heap.Fix(&workList, pos)
				}
			}
		}
	}
	return -1
}

func (c Coord) Heuristic(dst *Coord) int {
	// Manhattan distance, assuming min of 1 per traverse.
	return (dst.C - c.C) + (dst.R - c.R)
}

func (r RiskMap) Neighbors(pos Coord) []Coord {
	var coords []Coord
	up := Coord{pos.R - 1, pos.C}
	down := Coord{pos.R + 1, pos.C}
	left := Coord{pos.R, pos.C - 1}
	right := Coord{pos.R, pos.C + 1}
	for _, v := range []Coord{up, down, left, right} {
		if _, ok := r[v]; ok {
			coords = append(coords, v)
		}
	}
	return coords
}
