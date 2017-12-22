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
		score, garbage := process(l)
		fmt.Printf("Score for input with length %d: %d with %d removed garbage\n", len(l), score, garbage)
	}
}

func process(s string) (int, int) {
	totalScore := 0
	groupScore := 0
	removedGarbage := 0
	var inGarbage, escaped bool
	for _, c := range s {
		if escaped {
			escaped = false
			continue
		}
		if inGarbage {
			switch c {
			case '>':
				inGarbage = false
			case '!':
				escaped = true
			default:
				removedGarbage++
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
	return totalScore, removedGarbage
}
