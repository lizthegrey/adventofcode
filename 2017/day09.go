package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day09.input", "Relative file path to use as input.")

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes[:len(bytes)-1])
	for _, l := range strings.Split(contents, "\n") {
		fmt.Printf("Score for input with length %d: %d\n", len(l), score(l))
	}
}

func score(s string) int {
	totalScore := 0
	groupScore := 0
	var inGarbage, escaped bool
	for _, c := range s {
		if escaped {
			escaped = false
			continue
		}
		if inGarbage {
			if c == '>' {
				inGarbage = false
			}
			if c == '!' {
				escaped = true
			}
			continue
		}
		// We are in a non-garbage group.
		switch c {
		case '<':
			inGarbage = true
		case '{':
			groupScore++
		case '}':
			totalScore += groupScore
			groupScore--
		case ',':
			// Nothing needs to happen here.
		case '!':
			escaped = true
		default:
			// Nothing needs to happen here.
		}
	}
	return totalScore
}
