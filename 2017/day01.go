package main

import (
	"flag"
	"fmt"
	"strconv"
)

var input = flag.String("input", "1122", "The input to the problem.")
var partB = flag.Bool("partB", true, "Whether to use the Part B logic.")

func main() {
	flag.Parse()

	lastVal := -1  // Initialized to value that should never be a real digit.
	runningSum := 0
	firstVal := -1

	// digitsSeen := make([]int, len(*input))

	for i, c := range *input {
		curVal, err := strconv.Atoi(string(c))
		if err != nil {
			fmt.Printf("Couldn't parse: %v\n", err)
			return
		}
		if i == 0 {
			firstVal = curVal
		}
		if curVal == lastVal {
			runningSum += curVal
		}
		lastVal = curVal
	}
	if lastVal == firstVal {
		runningSum += lastVal
	}

	fmt.Printf("%d\n", runningSum)
}
