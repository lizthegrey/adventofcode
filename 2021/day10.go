package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day10.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	scoreA := 0
	scoreB := make([]int, 0)
outer:
	for _, line := range split {
		stack := make([]rune, 0)
		for _, c := range line {
			switch c {
			case '{':
				stack = append(stack, '}')
			case '(':
				stack = append(stack, ')')
			case '[':
				stack = append(stack, ']')
			case '<':
				stack = append(stack, '>')
			case '>':
				if len(stack) == 0 || stack[len(stack)-1] != c {
					scoreA += 25137
					continue outer
				}
				stack = stack[:len(stack)-1]
			case ']':
				if len(stack) == 0 || stack[len(stack)-1] != c {
					scoreA += 57
					continue outer
				}
				stack = stack[:len(stack)-1]
			case ')':
				if len(stack) == 0 || stack[len(stack)-1] != c {
					scoreA += 3
					continue outer
				}
				stack = stack[:len(stack)-1]
			case '}':
				if len(stack) == 0 || stack[len(stack)-1] != c {
					scoreA += 1197
					continue outer
				}
				stack = stack[:len(stack)-1]
			}
		}
		// Part B needs to be scored.
		intermediate := 0
		for i := len(stack) - 1; i >= 0; i-- {
			intermediate *= 5
			switch stack[i] {
			case ')':
				intermediate += 1
			case ']':
				intermediate += 2
			case '}':
				intermediate += 3
			case '>':
				intermediate += 4
			}
		}
		scoreB = append(scoreB, intermediate)
	}
	fmt.Println(scoreA)
	sort.Ints(scoreB)
	fmt.Println(scoreB[len(scoreB)/2])
}
