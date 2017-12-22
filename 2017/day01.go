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

	runningSum := 0
	digits := make([]int, len(*input))

	for i, c := range *input {
		curVal, err := strconv.Atoi(string(c))
		if err != nil {
			fmt.Printf("Couldn't parse: %v\n", err)
			return
		}
		digits[i] = curVal
	}

	for i, curVal := range digits {
		if *partB {
			if digits[(len(digits)/2+i)%len(digits)] == curVal {
				runningSum += curVal
			}
		} else {
			if i == 0 {
				if digits[len(digits)-1] == curVal {
					runningSum += curVal
				}
			} else if curVal == digits[i-1] {
				runningSum += curVal
			}
		}
	}

	fmt.Printf("%d\n", runningSum)
}
