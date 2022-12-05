package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day05.input", "Relative file path to use as input.")

type stacks [][]byte

func (s stacks) insertBottom(col int, val byte) stacks {
	if len(s) <= col {
		backing := make(stacks, col+1)
		copy(backing, s)
		s = backing
	}
	if s[col] == nil {
		s[col] = make([]byte, 0)
	}
	s[col] = append([]byte{val}, s[col]...)
	return s
}

func (s stacks) move(from, to, count int) {
	vals := s[from][len(s[from])-count : len(s[from])]
	s[from] = s[from][:len(s[from])-count]
	s[to] = append(s[to], vals...)
}

func (s stacks) printTop() {
	for i := 0; i < len(s); i++ {
		if len(s[i]) == 0 {
			fmt.Printf(" ")
			continue
		}
		fmt.Printf("%c", s[i][len(s[i])-1])
	}
	fmt.Println()
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var crates, cratesB stacks
	for _, s := range split[:len(split)-1] {
		if strings.Contains(s, "[") {
			for i := 1; i < len(s); i += 4 {
				if s[i] != ' ' {
					crates = crates.insertBottom(i/4, s[i])
					cratesB = cratesB.insertBottom(i/4, s[i])
				}
			}
		} else if len(s) > 5 && s[0:4] == "move" {
			parts := strings.Split(s, " ")
			count, _ := strconv.Atoi(parts[1])
			from, _ := strconv.Atoi(parts[3])
			to, _ := strconv.Atoi(parts[5])
			for i := 0; i < count; i++ {
				crates.move(from-1, to-1, 1)
			}
			cratesB.move(from-1, to-1, count)
		}
	}
	crates.printTop()
	cratesB.printTop()
}
