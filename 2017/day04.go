package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day04.input", "Relative file path to use for input.")
var partB = flag.Bool("partB", true, "Whether to use part B logic.")

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	valid := 0
outer:
	for _, line := range strings.Split(string(bytes), "\n") {
		if len(line) == 0 {
			break
		}
		fields := strings.Fields(line)
		seen := make(map[string]bool)
		for _, word := range fields {
			if seen[word] {
				continue outer
			}
			seen[word] = true
		}
		valid++
	}
	fmt.Printf("Found %d valid passphrases.\n", valid)
}
