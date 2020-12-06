package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day06.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]
	groupNum := 0
	groups := []map[rune]int{
		make(map[rune]int),
	}
	groupSize := []int{0}
	for _, s := range split {
		if s == "" {
			groups = append(groups, make(map[rune]int))
			groupSize = append(groupSize, 0)
			groupNum++
			continue
		}
		for _, r := range s {
			groups[groupNum][r]++
		}
		groupSize[groupNum]++
	}
	sum := 0
	all := 0
	for n, g := range groups {
		sum += len(g)
		for _, count := range g {
			if count == groupSize[n] {
				all++
			}
		}
	}
	fmt.Println(sum)
	fmt.Println(all)
}
