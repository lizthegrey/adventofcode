package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day02.input", "Relative file path to use as input.")

var max = map[string]int{
	"red":   12,
	"green": 13,
	"blue":  14,
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	sum := 0
	powers := 0
	for i, s := range split[:len(split)-1] {
		results := strings.Split(strings.Split(s, ": ")[1], "; ")
		possible := true
		seen := make(map[string]int)
		for _, r := range results {
			for _, v := range strings.Split(r, ", ") {
				parts := strings.Split(v, " ")
				count, err := strconv.Atoi(parts[0])
				if err != nil {
					log.Fatalf("%s failed: %v", parts[0], err)
				}
				colour := parts[1]
				if max[colour] < count {
					possible = false
				}
				if seen[colour] < count {
					seen[colour] = count
				}
			}
		}
		if possible {
			sum += i + 1
		}
		power := 1
		for _, v := range seen {
			power *= v
		}
		powers += power
	}
	fmt.Println(sum)
	fmt.Println(powers)
}
