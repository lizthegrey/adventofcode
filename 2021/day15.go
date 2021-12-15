package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
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
	workList := map[Coord]bool{
		*src: true,
	}
	history := make(map[Coord]Coord)

	for len(workList) != 0 {
		// Pop the current node off the worklist.
		currentScore := math.MaxInt
		var current Coord
		for v := range workList {
			if score, ok := fScore[v]; ok {
				if score < currentScore {
					current = v
					currentScore = score
				}
			}
		}
		delete(workList, current)

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
				if !workList[n] {
					workList[n] = true
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
