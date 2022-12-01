package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
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
	var highest int
	var current int
	for _, s := range split {
		if s == "" {
			if highest < current {
				highest = current
			}
			current = 0
		}
		i, _ := strconv.Atoi(s)
		current += i
	}
	if highest < current {
		highest = current
	}
	fmt.Println(highest)

	// part B
	var found []int
	current = 0
	for _, s := range split {
		if s == "" {
			found = append(found, current)
			current = 0
		}
		i, _ := strconv.Atoi(s)
		current += i
	}
	found = append(found, current)
	sort.Ints(found)
	var sum int
	for i := 1; i <= 3; i++ {
		sum += found[len(found)-i]
	}
	fmt.Println(sum)
}
