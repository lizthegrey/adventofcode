package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day04.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var countA, countB int
	for _, s := range split[:len(split)-1] {
		pairs := strings.Split(s, ",")
		first := strings.Split(pairs[0], "-")
		second := strings.Split(pairs[1], "-")
		firstStart, _ := strconv.Atoi(first[0])
		firstEnd, _ := strconv.Atoi(first[1])
		secondStart, _ := strconv.Atoi(second[0])
		secondEnd, _ := strconv.Atoi(second[1])
		if firstStart >= secondStart && firstEnd <= secondEnd {
			countA++
			countB++
			continue
		}
		if firstStart <= secondStart && firstEnd >= secondEnd {
			countA++
			countB++
			continue
		}
		// Not a complete overlap, but check for partial overlap
		if firstStart >= secondStart && firstStart <= secondEnd {
			countB++
			continue
		}
		if secondStart >= firstStart && secondStart <= firstEnd {
			countB++
			continue
		}
		if firstEnd >= secondStart && firstEnd <= secondEnd {
			countB++
			continue
		}
		if secondEnd >= firstStart && secondEnd <= firstEnd {
			countB++
			continue
		}
	}
	fmt.Println(countA)
	fmt.Println(countB)
}
