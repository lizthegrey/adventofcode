package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
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
	split = split[:len(split)-1]
	var matchesA, matchesB int
	for _, s := range split {
		parsed := strings.Split(s, " ")
		if len(parsed) != 3 {
			fmt.Printf("Encountered bad line: %s\n", split)
		}

		lenSpec := strings.Split(parsed[0], "-") // "1-3"
		if len(lenSpec) != 2 {
			fmt.Printf("Encountered bad lenSpec: %s\n", parsed[0])
		}
		char := rune(parsed[1][0]) // "n:"
		password := parsed[2]      // "xxxxxxx"

		min, err := strconv.Atoi(lenSpec[0])
		max, err := strconv.Atoi(lenSpec[1])
		if err != nil {
			fmt.Printf("Failed to parse %s\n", parsed[0])
			break
		}
		chars := make(map[rune]int)
		for _, c := range password {
			chars[c]++
		}
		if chars[char] >= min && chars[char] <= max {
			matchesA++
		}
		if (rune(password[min-1]) == char) != (rune(password[max-1]) == char) {
			matchesB++
		}
	}
	fmt.Println(matchesA)
	fmt.Println(matchesB)
}
