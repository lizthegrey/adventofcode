package main

import (
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
)

var debug = flag.Bool("debug", false, "Print debug info as we go along.")
var partB = flag.Bool("partB", false, "Use part B logic.")

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

	var phaseList []int
	if !*partB {
		phaseList = []int{0, 1, 2, 3, 4}
	} else {
		phaseList = []int{5, 6, 7, 8, 9}
	}
	highestOutput := -1
	permutations := permute(phaseList)
	if len(permutations) != 120 {
		fmt.Printf("Failed to get right permutations: %d\n", len(permutations))
	}
	for _, phases := range permutations {
		inputVal := 0
		inputs := make([]chan int, len(phases))
		outputs := make([]chan int, len(phases))
		dones := make([]chan bool, len(phases))
		for i, p := range phases {
			workingTape := make(intcode.Tape, len(tape))
			copy(workingTape, tape)
			input := make(chan int)
			output, done := workingTape.Process(input)
			input <- p
			if !*partB || i == 0 {
				input <- inputVal
			}
			if !*partB {
				inputVal = <-output
			} else {
				inputs[i] = input
				outputs[i] = output
				dones[i] = done
			}
		}
		if !*partB && (inputVal > highestOutput) {
			highestOutput = inputVal
		}

		if *partB {
			f := make(chan bool)
			for i := range inputs {
				go func(idx int) {
					oIdx := idx - 1
					if idx == 0 {
						oIdx = len(inputs) - 1
					}
					for {
						val := <-outputs[oIdx]
						select {
						case <-dones[idx]:
							if highestOutput < val {
								highestOutput = val
							}
							f <- true
							return
						default:
							fmt.Printf("Feeding %d from %d to %d\n", val, oIdx, idx)
							inputs[idx] <- val
						}
					}
				}(i)
			}
			<-f
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
