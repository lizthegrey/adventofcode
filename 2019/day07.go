package main

import (
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
)

var debug = flag.Bool("debug", false, "Print debug info as we go along.")

func main() {
	flag.Parse()
	tape := intcode.ReadInput()
	if tape == nil {
		fmt.Println("Failed to parse input.")
		return
	}
	if *debug {
		for _, n := range tape {
			fmt.Printf("%d,", n)
		}
		fmt.Println()
	}

	phaseList := []int{0, 1, 2, 3, 4}
	highestOutput := -1
	permutations := permute(phaseList)
	if len(permutations) != 120 {
		fmt.Printf("Failed to get right permutations: %d\n", len(permutations))
	}
	for _, phases := range permutations {
		input := 0
		for _, p := range phases {
			workingTape := make(intcode.Tape, len(tape))
			copy(workingTape, tape)
			_, input = workingTape.Process([]int{p, input})
		}
		if input > highestOutput {
			highestOutput = input
		}
	}
	fmt.Println(highestOutput)
}

func permute(in []int) [][]int {
	ret := make([][]int, 0)
	if len(in) == 1 {
		return [][]int{{in[0]}}
	}
	for i, v := range in {
		// Put v at the front, then use all the permutations of the rest.
		rest := make([]int, 0)
		rest = append(rest, in[0:i]...)
		rest = append(rest, in[i+1:]...)
		for _, tail := range permute(rest) {
			candidate := make([]int, 1)
			candidate[0] = v
			candidate = append(candidate, tail...)
			ret = append(ret, candidate)
		}
	}
	return ret
}
