package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day01.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]
	seen := make([]int, len(split))
	for i, s := range split {
		n, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("Failed to parse %s\n", s)
			break
		}
		seen[i] = n
	}
partA:
	for i, n := range seen {
		for j, m := range seen {
			if i >= j {
				continue
			}
			if n+m == 2020 {
				fmt.Println(n * m)
				break partA
			}
		}
	}
partB:
	for i, n := range seen {
		for j, m := range seen {
			for k, o := range seen {
				if i >= j || j >= k {
					continue
				}
				if n+m+o == 2020 {
					fmt.Println(n * m * o)
					break partB
				}
			}
		}
	}
}
