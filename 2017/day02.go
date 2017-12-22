package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day02.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", true, "Whether to use part B logic.")

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes)
	result := 0
	for _, line := range strings.Split(contents, "\n") {
		if len(line) == 0 {
			break
		}
		largest := -1
		smallest := math.MaxInt16
		fields := strings.Fields(line)
		seen := make([]int, len(fields))
		for i, num := range fields {
			n, err := strconv.Atoi(num)
			seen[i] = n
			if err != nil {
				fmt.Printf("Could not parse %s because %v.\n", num, err)
			}
			if n < smallest {
				smallest = n
			}
			if n > largest {
				largest = n
			}
			if *partB {
				for j := 0; j < i; j++ {
					dividend, divisor := -1, -1
					if seen[j] < n {
						dividend = n
						divisor = seen[j]
					} else {
						dividend = seen[j]
						divisor = n
					}
					if dividend%divisor == 0 {
						result += dividend / divisor
					}
				}
			}
		}
		if !*partB {
			result += largest - smallest
		}
	}
	fmt.Printf("Result is %d\n", result)
}
