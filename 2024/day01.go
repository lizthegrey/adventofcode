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

	var listA, listB []int
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, " ")
		first := true
		for _, p := range parts {
			if p == "" {
				// hacky but this works instead of writing a regex matcher.
				continue
			}
			n, _ := strconv.Atoi(p)
			if !first {
				listB = append(listB, n)
			} else {
				listA = append(listA, n)
			}
			first = false
		}
	}

	// part A
	sort.Ints(listA)
	sort.Ints(listB)
	var distance int
	for i := range listA {
		diff := listB[i] - listA[i]
		if diff < 0 {
			diff *= -1
		}
		distance += diff
	}
	fmt.Println(distance)

	// part B
	frequencies := make(map[int]int)
	for _, b := range listB {
		frequencies[b] += 1
	}
	var similarity int
	for _, a := range listA {
		similarity += a * frequencies[a]
	}
	fmt.Println(similarity)
}
