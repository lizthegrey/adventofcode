package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day23.input", "Relative file path to use as input.")

type Coord struct {
	R, C int
}

type Grid map[Coord]rune

type Hist struct {
	Coord
	Distance int
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	grid := make(Grid)
	maxC := len(split[0]) - 1
	var maxR int
	for r, row := range split[:len(split)-1] {
		for c, v := range row {
			grid[Coord{r, c}] = v
		}
		maxR = r
	}
	start := Coord{0, 1}
	target := Coord{maxR, maxC - 1}
	fmt.Println(grid.Visit([]Hist{{start, 0}}, target, true))
	fmt.Println(grid.Visit([]Hist{{start, 0}}, target, false))
}

// DFS/recursion should be fine here for a max depth of a few hundred..
func (g Grid) Visit(hist []Hist, target Coord, enforceSlopes bool) int {
	histLen := len(hist)
	src := hist[histLen-1]
	if src.Coord == target {
		var sum int
		for _, v := range hist {
			sum += v.Distance
		}
		return sum
	}
	longest := -1
	neighs := g.Neighbours(hist, enforceSlopes)
	for _, neigh := range neighs {
		hist = hist[:histLen]
		hist = append(hist, neigh)
		if pathLen := g.Visit(hist, target, enforceSlopes); pathLen > longest {
			longest = pathLen
		}
	}
	return longest
}

func (c Coord) Adjacent() []Coord {
	return []Coord{{c.R - 1, c.C}, {c.R, c.C - 1}, {c.R, c.C + 1}, {c.R + 1, c.C}}
}

func (g Grid) Neighbours(hist []Hist, enforceSlopes bool) []Hist {
	var ret []Hist
	src := hist[len(hist)-1].Coord
outer:
	for _, neigh := range src.Adjacent() {
		// Enforce one-ways.
		if enforceSlopes {
			switch g[src] {
			case '^':
				if neigh.R >= src.R {
					continue
				}
			case 'v':
				if neigh.R <= src.R {
					continue
				}
			case '<':
				if neigh.C >= src.C {
					continue
				}
			case '>':
				if neigh.C <= src.C {
					continue
				}
			}
		}

		dist := 1
		prev := src
	inner:
		for {
			switch g[neigh] {
			case 0:
				// We've gone off the board.
				continue outer
			case '#':
				continue outer
			case '.':
				// Ordinary case.
			case '^':
				if neigh.R > src.R && enforceSlopes {
					continue outer
				}
			case 'v':
				if neigh.R < src.R && enforceSlopes {
					continue outer
				}
			case '<':
				if neigh.C > src.C && enforceSlopes {
					continue outer
				}
			case '>':
				if neigh.C < src.C && enforceSlopes {
					continue outer
				}
			}
			neighneigh := neigh.Adjacent()
			var count int
			var candidate Coord
			for _, nn := range neighneigh {
				if nn == prev {
					// Disallow backtracking
					continue
				}
				if g[nn] != 0 && g[nn] != '#' {
					count++
					candidate = nn
				}
			}
			switch count {
			case 1:
				dist++
				prev = neigh
				neigh = candidate
			default:
				// We're at a junction or dead end.
				break inner
			}
		}

		for i := len(hist) - 1; i >= 0; i-- {
			// Optimisation: hist is most likely to backtrack our most recent step
			if hist[i].Coord == neigh {
				continue outer
			}
		}
		ret = append(ret, Hist{neigh, dist})
	}
	return ret
}
