package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day06.input", "Relative file path to use as input.")
var debug = flag.Bool("debug", false, "Whether to print debug output.")

var com string = "COM"
var me string = "YOU"
var santa string = "SAN"

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	// Inwards: outer -> center
	parents := make(map[string]string)

	for _, s := range split {
		if s == "" {
			continue
		}
		line := strings.Split(s, ")")
		if len(line) != 2 {
			fmt.Printf("Failed to parse line %s\n", s)
			return
		}
		center := line[0]
		outer := line[1]
		parents[outer] = center
	}

	distances := make(map[string]int)
	distances[com] = 0
	sum := 0
	for k := range parents {
		sum += computeDistance(k, distances, parents)
	}
	fmt.Printf("Part A: %d\n", sum)

	if parents[me] == "" || parents[santa] == "" {
		fmt.Println("Couldn't find YOU or SAN.")
		return
	}

	myParents := make([]string, 0)
	santaParents := make([]string, 0)
	for loc := me; loc != com; loc = parents[loc] {
		myParents = append(myParents, loc)
	}
	for loc := santa; loc != com; loc = parents[loc] {
		santaParents = append(santaParents, loc)
	}

	// Step backwards through the lists.
	for i := 1; ; i++ {
		s := santaParents[len(santaParents)-i]
		m := myParents[len(myParents)-i]
		if *debug {
			fmt.Printf("Step %d: %s, %s\n", i, s, m)
		}
		if s != m {
			fmt.Printf("Part B: %d\n", len(myParents)+len(santaParents)-(2*i))
			break
		}
	}
}

func computeDistance(object string, distances map[string]int, parents map[string]string) int {
	if _, ok := distances[object]; !ok {
		distances[object] = computeDistance(parents[object], distances, parents) + 1
	}
	return distances[object]
}
