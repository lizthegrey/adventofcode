package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day15.input", "Relative file path to use as input.")
var debug = flag.Bool("debug", false, "Whether to print debug output along the way.")
var partB = flag.Bool("partB", false, "Whether to use part B logic.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]
	split = strings.Split(split[0], ",")
	seen := make(map[int]int)
	prev := -1
	for i, s := range split {
		n, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("Failed to parse %s\n", s)
			break
		}
		seen[n] = i
		prev = n
	}
	max := 2020
	if *partB {
		max = 30000000
	}
	for i := len(split); i < max; i++ {
		var val int
		if prevPos, seen := seen[prev]; seen {
			val = (i - 1) - prevPos
		} else {
			val = 0
		}
		if *debug {
			fmt.Println(val)
		}
		if prev >= 0 {
			seen[prev] = i - 1
		}
		prev = val
	}
	fmt.Println(prev)
}
