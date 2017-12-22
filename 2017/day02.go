package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "day02.input", "Relative file path to use as input.")

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
		for _, num := range strings.Fields(line) {
			n, err := strconv.Atoi(num)
			if err != nil {
				fmt.Printf("Could not parse %s because %v.\n", num, err)
			}
			if n < smallest {
				smallest = n
			}
			if n > largest {
				largest = n
			}
		}
		result += largest - smallest
	}
	fmt.Printf("Result is %d\n", result)
}
