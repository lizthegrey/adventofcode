package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day01.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	// part A
	last := -1
	increments := 0
	for _, s := range split {
		i, _ := strconv.Atoi(s)
		if last != -1 && i > last {
			increments++
		}
		last = i
	}
	fmt.Println(increments)

	// part B
	increments = 0
	// Keep a sliding window of the last 3 numbers seen.
	lastNums := []int{0, 0, 0}
	lastSum := 0
	for idx, s := range split {
		i, _ := strconv.Atoi(s)
		// Remove the oldest number and add ourselves.
		sum := lastSum + i - lastNums[0]
		// Roll the oldest number off.
		lastNums = append(lastNums[1:], i)

		if idx < 3 {
			lastSum = sum
			continue
		}
		if lastSum < sum {
			increments++
		}
		lastSum = sum
	}
	fmt.Println(increments)
}
