package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day13.input", "Relative file path to use as input.")
var debug = flag.Bool("debug", false, "Whether to print debug output along the way.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]
	startTime, err := strconv.Atoi(split[0])
	if err != nil {
		fmt.Printf("Failed to parse %s\n", split[0])
		return
	}
	lines := strings.Split(split[1], ",")
	soonest := startTime
	soonestNo := -1
	ids := make(map[int]int)
	for i, l := range lines {
		if l == "x" {
			continue
		}
		lineNo, err := strconv.Atoi(l)
		if err != nil {
			fmt.Printf("Failed to parse %s\n", l)
			return
		}
		ids[lineNo] = i
		minToWait := lineNo - (startTime % lineNo)
		if minToWait < soonest {
			soonest = minToWait
			soonestNo = lineNo
		}
	}
	fmt.Println(soonest * soonestNo)

	// Part B.
	// We are looking for the following relation to be satisfied:
	// ex t st t + v % ids[k] === 0 for each k,v
	minValue := 0
	runningProduct := 1
	for k, v := range ids {
		for (minValue+v)%k != 0 {
			minValue += runningProduct
		}
		runningProduct *= k
		if *debug {
			fmt.Printf("t + %d === 0 mod %d\n", v, k)
			fmt.Printf("Sum so far: %d, product so far: %d\n", minValue, runningProduct)
		}
	}
	fmt.Println(minValue)
}
