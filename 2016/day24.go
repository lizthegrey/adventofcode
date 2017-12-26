package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day24.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", true, "Return to 0 after all other nodes visited.")

type Coord struct {
	Row, Col int
}

type Pair struct {
	Src, Dst int
}

type CoordList []Coord

type Maze map[Coord]bool

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes)
	mazeRaw := strings.Split(contents[:len(contents)-1], "\n")

	pointsOfInterest := make(map[int]Coord)
	maze := make(Maze, len(mazeRaw))
	for r, l := range mazeRaw {
		for c, point := range l {
			switch point {
			case '#':
				// Don't insert into the passable map.
			case '.':
				// Mark the node passable.
				maze[Coord{r, c}] = true
			default:
				poi, err := strconv.Atoi(string(point))
				if err != nil {
					fmt.Printf("Failed to parse POI: %c", point)
				}
				pointsOfInterest[poi] = Coord{r, c}

				// Mark the node passable.
				maze[Coord{r, c}] = true
			}
		}
	}

	// Find shortest paths between each pair of POIs. (8*7/2 = 28 different routes)
	memo := make(map[Pair]int)
	for k1, v1 := range pointsOfInterest {
		search := []Coord{v1}
		visited := map[Coord]int{v1: 0}

		for {
			search = maze.expandRing(visited, search)

			missing := false
			for k2, v2 := range pointsOfInterest {
				if k2 >= k1 {
					continue
				}
				if dist := visited[v2]; dist == 0 {
					missing = true
				} else {
					memo[Pair{k1, k2}] = dist
					memo[Pair{k2, k1}] = dist
				}
			}
			if !missing {
				break
			}
		}
	}

	// Generate permutations of POIs to visit, starting from 0 (factorial(7) ~~ 5000)
	poiIndices := make([]int, len(pointsOfInterest)-1)
	i := 0
	for k, _ := range pointsOfInterest {
		if k != 0 {
			poiIndices[i] = k
			i++
		}
	}
	permutes := listPermutations(poiIndices)

	// Evaluate total route length for each possibility and emit the shortest sum.
	shortest := math.MaxUint16
	for _, p := range permutes {
		loc := 0
		dist := 0
		for _, v := range p {
			dist += memo[Pair{loc, v}]
			loc = v
		}
		if *partB {
			dist += memo[Pair{loc, 0}]
		}
		if dist < shortest {
			shortest = dist
		}
	}
	fmt.Println(shortest)
}

// Expand the ring of what we've visited outwards until we find the coord we're looking for.
// Maximum complexity of this is O(R*C*STARTS) since each passable square visited once per start.
func (m Maze) expandRing(visited map[Coord]int, toSearch []Coord) []Coord {
	ret := make([]Coord, 0)
	for _, start := range toSearch {
		depth := visited[start]
		up := Coord{start.Row - 1, start.Col}
		right := Coord{start.Row, start.Col + 1}
		down := Coord{start.Row + 1, start.Col}
		left := Coord{start.Row, start.Col - 1}

		neighs := []Coord{up, right, down, left}
		for _, n := range neighs {
			if _, found := visited[n]; found {
				continue
			} else {
				if m[n] {
					visited[n] = depth + 1
					ret = append(ret, n)
				}
			}
		}
	}
	return ret
}

func listPermutations(places []int) [][]int {
	ret := make([][]int, 0)

	if len(places) == 1 {
		return [][]int{[]int{places[0]}}
	}

	for _, p := range places {
		temp := make([]int, len(places)-1)
		i := 0
		for _, v := range places {
			if v != p {
				temp[i] = v
				i++
			}
		}
		lps := listPermutations(temp)
		for i, v := range lps {
			here := []int{p}
			lps[i] = append(here, v...)
		}

		ret = append(ret, lps...)
	}

	return ret
}
