package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day12.input", "Relative file path to use as input.")

type Program struct {
	Neighbors []int
}
type Progs map[int]Program

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes)
	lines := strings.Split(contents, "\n")
	programs := make(Progs)

	for _, l := range lines {
		if len(l) == 0 {
			break
		}
		parts := strings.Split(l, " <-> ")

		n, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Printf("Failed to parse program number in '%s'\n", l)
			return
		}

		neighbors := strings.Split(parts[1], ", ")
		neighs := make([]int, len(neighbors))
		for i, v := range neighbors {
			neigh, err := strconv.Atoi(v)
			if err != nil {
				fmt.Printf("Failed to parse neighbor number in '%s'\n", l)
				return
			}
			neighs[i] = neigh
		}
		programs[n] = Program{neighs}
	}
	visited := make(map[int]bool)
	programs.DeepTraverse(0, visited)
	fmt.Println(len(visited))
}

func (p Progs) DeepTraverse(i int, visited map[int]bool) {
	prog := p[i]
	if visited[i] {
		return
	}
	visited[i] = true
	for _, n := range prog.Neighbors {
		p.DeepTraverse(n, visited)
	}
}
