package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day07.input", "Relative file path to use as input.")

type Coord struct {
	R, C int
}
type Memo map[Coord]uint64

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)[:len(bytes)-1]

	var start Coord
	splitters := make(map[Coord]bool)
	beam := make(map[Coord]bool)
	var maxR int
	for r, s := range strings.Split(contents, "\n") {
		for c, v := range s {
			loc := Coord{r, c}
			switch v {
			case 'S':
				start = loc
			case '.':
				// pass
			case '^':
				splitters[loc] = true
			}
		}
		maxR = r
	}

	worklist := []Coord{start}
	children := make(map[Coord][]Coord)
	for len(worklist) > 0 {
		cur := worklist[0]
		worklist = worklist[1:]
		if len(children[cur]) > 0 {
			// This node has been processed already.
			continue
		}

		if !splitters[cur] {
			// Initial setup conditions starting from start
			beam[cur] = true
			worklist = append(worklist, Coord{cur.R + 1, cur.C})
			continue
		}

		if len(children[start]) == 0 {
			// This is the initial setup phase completion.
			children[start] = append(children[start], cur)
		}

		for _, branch := range []Coord{{cur.R, cur.C - 1}, {cur.R, cur.C + 1}} {
			for branch.R <= maxR {
				if splitters[branch] {
					children[cur] = append(children[cur], branch)
					worklist = append(worklist, branch)
					break
				}
				beam[branch] = true
				branch.R++
			}
		}
	}

	var countA int
	counted := make(map[Coord]bool)
	for _, dsts := range children {
		for _, dst := range dsts {
			if !counted[dst] {
				countA++
				counted[dst] = true
			}
		}
	}
	fmt.Println(countA)

	m := make(Memo)
	fmt.Println(m.compute(children, children[start][0]))
}

func (m Memo) compute(children map[Coord][]Coord, c Coord) uint64 {
	if val, ok := m[c]; ok {
		return val
	}
	ways := uint64(2 - len(children[c]))
	for _, child := range children[c] {
		ways += m.compute(children, child)
	}
	m[c] = ways
	return ways
}
