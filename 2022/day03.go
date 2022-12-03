package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day03.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	// part A
	var priorities int
	for _, s := range split {
		if s == "" {
			break
		}
		inFirst := make(map[rune]bool)
		first := s[0 : len(s)/2]
		second := s[len(s)/2 : len(s)]
		for _, k := range first {
			inFirst[k] = true
		}
		for _, k := range second {
			if inFirst[k] {
				if k >= 'a' && k <= 'z' {
					priorities += int(k-'a') + 1
				} else {
					priorities += int(k-'A') + 27
				}
				break
			}
		}
	}
	fmt.Println(priorities)

	// part B
	priorities = 0
	for i := 0; i < len(split)-1; i += 3 {
		found := make(map[rune]int)
		for j := 0; j < 3; j++ {
			local := make(map[rune]bool)
			for _, k := range split[i+j] {
				if local[k] {
					continue
				}
				local[k] = true
				found[k]++
			}
		}
		for k, v := range found {
			if v == 3 {
				if k >= 'a' && k <= 'z' {
					priorities += int(k-'a') + 1
				} else {
					priorities += int(k-'A') + 27
				}
				break
			}
		}
	}
	fmt.Println(priorities)
}
