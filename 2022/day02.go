package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day02.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	// part A
	var score int
	for _, s := range split {
		if s == "" {
			break
		}
		vals := strings.Split(s, " ")
		opp := vals[0]
		mine := vals[1]
		switch mine {
		case "X":
			score += 1
		case "Y":
			score += 2
		case "Z":
			score += 3
		}
		switch opp {
		case "A":
			switch mine {
			case "X":
				score += 3
			case "Y":
				score += 6
			case "Z":
				score += 0
			}
		case "B":
			switch mine {
			case "X":
				score += 0
			case "Y":
				score += 3
			case "Z":
				score += 6
			}
		case "C":
			switch mine {
			case "X":
				score += 6
			case "Y":
				score += 0
			case "Z":
				score += 3
			}
		}
	}
	fmt.Println(score)

	// part B
	score = 0
	for _, s := range split {
		if s == "" {
			break
		}
		vals := strings.Split(s, " ")
		opp := vals[0]
		outcome := vals[1]
		switch outcome {
		case "X":
			score += 0
		case "Y":
			score += 3
		case "Z":
			score += 6
		}
		switch opp {
		case "A":
			switch outcome {
			case "X":
				score += 3
			case "Y":
				score += 1
			case "Z":
				score += 2
			}
		case "B":
			switch outcome {
			case "X":
				score += 1
			case "Y":
				score += 2
			case "Z":
				score += 3
			}
		case "C":
			switch outcome {
			case "X":
				score += 2
			case "Y":
				score += 3
			case "Z":
				score += 1
			}
		}
	}
	fmt.Println(score)
}
