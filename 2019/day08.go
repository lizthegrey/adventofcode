package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
)

var inputFile = flag.String("inputFile", "inputs/day08.input", "Relative file path to use as input.")

const cols int = 25
const rows int = 6

type Sheet [rows][cols]int

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes[:len(bytes)-1])
	layers := make([]Sheet, 0)
	layer := -1
	for i, s := range contents {
		n, err := strconv.Atoi(string(s))
		if err != nil {
			fmt.Printf("Failed to parse %s\n", s)
			break
		}
		if i%(cols*rows) == 0 {
			layer++
			layers = append(layers, Sheet{})
		}
		idx := i % (cols * rows)
		layers[layer][idx/cols][idx%cols] = n
	}
	fewestZeroes := cols * rows
	ret := -1
	for _, layer := range layers {
		seen := make(map[int]int)
		for r := range layer {
			for c := range layer[r] {
				seen[layer[r][c]]++
			}
		}
		if fewestZeroes > seen[0] {
			fewestZeroes = seen[0]
			ret = seen[1] * seen[2]
		}
	}
	fmt.Printf("Part A: %d\n", ret)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
		outer:
			for _, sheet := range layers {
				switch sheet[r][c] {
				case 0:
					fmt.Printf(" ")
					break outer
				case 1:
					fmt.Printf("x")
					break outer
				case 2:
					continue outer
				}
			}
		}
		fmt.Println()
	}
}
