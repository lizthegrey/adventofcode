package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
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
			if *partB {
				word = normalize(word)
			}
			if seen[word] {
				continue outer
			}
			seen[word] = true
		}
		valid++
	}
	fmt.Printf("Found %d valid passphrases.\n", valid)
}

type ByRune []rune

func (a ByRune) Len() int           { return len(a) }
func (a ByRune) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRune) Less(i, j int) bool { return a[i] < a[j] }

func normalize(s string) string {
	var rs []rune = []rune(s)
	sort.Sort(ByRune(rs))
	return string(rs)
}
