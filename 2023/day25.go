package main

import (
	"cmp"
	"flag"
	"fmt"
	"io/ioutil"
	"slices"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day25.input", "Relative file path to use as input.")

type Pair struct {
	Src, Dst string
}

type WeightedPair struct {
	Pair
	Weight int
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	edges := make(map[Pair]bool)
	vertices := make(map[string]bool)
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, ": ")
		src := parts[0]
		vertices[src] = true
		for _, dst := range strings.Split(parts[1], " ") {
			edges[Pair{src, dst}] = true
			edges[Pair{dst, src}] = true
			vertices[dst] = true
		}
	}

	exempt := make(map[string]bool)
	for i := 0; i <= 3; i++ {
		// This could be made a lot faster by updating an existing FW after edge removal. But shrug.
		dist, prev := FW(edges, vertices)
		if i == 3 {
			for n := range vertices {
				var reachable, unreachable int
				for m := range vertices {
					if _, ok := dist[Pair{n, m}]; ok {
						reachable++
					} else {
						unreachable++
					}
				}
				fmt.Println(reachable * unreachable)
				return
			}
		}

		weights := MostUsed(prev)
		for k := len(weights) - 1; k > 0; k-- {
			toClip := weights[k]
			edge := toClip.Pair
			if exempt[edge.Src] || exempt[edge.Dst] {
				continue
			}

			// fmt.Printf("Clipping %v\n", edge)
			exempt[edge.Src] = true
			exempt[edge.Dst] = true
			delete(edges, edge)
			delete(edges, Pair{edge.Dst, edge.Src})
			break
		}
	}
}

// Floyd-Warshall with path reconstruction.
func FW(edges map[Pair]bool, vertices map[string]bool) (map[Pair]int, map[Pair]string) {
	dist := make(map[Pair]int)
	prev := make(map[Pair]string)
	for e := range edges {
		dist[e] = 1
		prev[e] = e.Src
	}
	for v := range vertices {
		self := Pair{v, v}
		dist[self] = 0
		prev[self] = v
	}
	var n int
	for k := range vertices {
		n++
		for i := range vertices {
			for j := range vertices {
				ij := Pair{i, j}
				ik := Pair{i, k}
				kj := Pair{k, j}
				best, ok := dist[ij]
				a, okA := dist[ik]
				b, okB := dist[kj]
				if okA && okB && (!ok || a+b < best) {
					dist[ij] = a + b
					prev[ij] = prev[kj]
				}
			}
		}
	}
	return dist, prev
}

func MostUsed(prev map[Pair]string) []WeightedPair {
	mostUsed := make(map[Pair]int)
	for e, next := range prev {
		mostUsed[Pair{e.Dst, next}]++
		mostUsed[Pair{next, e.Dst}]++
	}
	var weights []WeightedPair
	for k, v := range mostUsed {
		weights = append(weights, WeightedPair{k, v})
	}
	slices.SortFunc(weights, func(a, b WeightedPair) int {
		return cmp.Compare(a.Weight, b.Weight)
	})
	return weights
}
