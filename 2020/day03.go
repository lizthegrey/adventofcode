package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day03.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]
	trees := make([][]bool, len(split))
	for i, s := range split {
		trees[i] = make([]bool, len(s))
		for j, c := range s {
			trees[i][j] = (c == '#')
		}
		if err != nil {
			fmt.Printf("Failed to parse %s\n", s)
			break
		}
	}

	fmt.Println(checkSlope(1, 3, trees))
	a := checkSlope(1, 3, trees)
	b := checkSlope(1, 1, trees)
	c := checkSlope(1, 5, trees)
	d := checkSlope(1, 7, trees)
	e := checkSlope(2, 1, trees)
	fmt.Printf("%d*%d*%d*%d*%d = %d\n", a, b, c, d, e, a*b*c*d*e)
}

func checkSlope(down, right int, trees [][]bool) int {
	hit := 0
	for time := 0; time*down < len(trees); time++ {
		traversed := time * down
		column := (time * right) % len(trees[traversed])
		if trees[traversed][column] {
			hit++
		}
	}
	return hit
}
