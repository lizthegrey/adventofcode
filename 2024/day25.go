package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day25.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var keys, locks [][5]int
	for i := 0; i < len(split[:len(split)-1]); i += 8 {
		var state [5]int
		for j := range 7 {
			for c, v := range split[i+j] {
				if v == '#' {
					state[c]++
				}
			}
		}
		if split[i][0] == '#' {
			locks = append(locks, state)
		} else {
			keys = append(keys, state)
		}
	}

	var matches int
	for _, k := range keys {
	outer:
		for _, l := range locks {
			for i := range 5 {
				if k[i]+l[i] > 7 {
					continue outer
				}
			}
			matches++
		}
	}
	fmt.Println(matches)
}
