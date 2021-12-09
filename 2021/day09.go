package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day09.input", "Relative file path to use as input.")

type HeightArray [][]int
type Coord struct {
	r, c int
}
type Basins map[Coord]Coord

func (h HeightArray) Grow(worklist []Coord, basins Basins) []Coord {
	newWorklist := make([]Coord, 0)
	for _, w := range worklist {
		myBasin := basins[w]
		for _, n := range h.Neighbors(w) {
			if _, ok := basins[n]; ok {
				// Already part of a basin, skip
				continue
			}
			if h[n.r][n.c] == 9 {
				// High points cannot be part of a basin
				continue
			}
			basins[n] = myBasin
			newWorklist = append(newWorklist, n)
		}
	}
	return newWorklist
}

func (h HeightArray) Neighbors(coord Coord) []Coord {
	ret := make([]Coord, 0)
	if coord.r > 0 {
		ret = append(ret, Coord{coord.r - 1, coord.c})
	}
	if coord.r < len(h)-1 {
		ret = append(ret, Coord{coord.r + 1, coord.c})
	}
	if coord.c > 0 {
		ret = append(ret, Coord{coord.r, coord.c - 1})
	}
	if coord.c < len(h[coord.r])-1 {
		ret = append(ret, Coord{coord.r, coord.c + 1})
	}
	return ret
}

func (h HeightArray) Lowest(coord Coord) bool {
	value := h[coord.r][coord.c]
	for _, n := range h.Neighbors(coord) {
		if h[n.r][n.c] <= value {
			return false
		}
	}
	return true
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	heights := make(HeightArray, 0)
	for _, s := range split {
		row := make([]int, 0)
		for _, v := range s {
			height := int(v - '0')
			row = append(row, height)
		}
		heights = append(heights, row)
	}

	worklist := make([]Coord, 0)
	basins := make(Basins)
	sum := 0
	for r, row := range heights {
		for c, height := range row {
			loc := Coord{r, c}
			if heights.Lowest(loc) {
				sum += height + 1
				basins[loc] = loc
				worklist = append(worklist, loc)
			}
		}
	}
	fmt.Println(sum)

	for len(worklist) != 0 {
		worklist = heights.Grow(worklist, basins)
	}
	basinSizes := make(map[Coord]int)
	for _, source := range basins {
		basinSizes[source]++
	}
	sizes := make([]int, 0)
	for _, v := range basinSizes {
		sizes = append(sizes, v)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
	fmt.Println(sizes[0] * sizes[1] * sizes[2])
}
