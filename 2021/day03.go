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
	split = split[:len(split)-1]

	var ones [12]int
	lines := 0
	data := make([][12]bool, 0)
	for _, s := range split {
		var line [12]bool
		lines++
		for pos, c := range s {
			if c == '1' {
				line[pos] = true
				ones[pos]++
			}
		}
		data = append(data, line)
	}
	var epsilon, gamma int
	for pos, count := range ones {
		if count*2 > lines {
			epsilon += 1 << (11 - pos)
		} else {
			gamma += 1 << (11 - pos)
		}
	}
	fmt.Println(epsilon * gamma)

	oxygen := iterate(data, func(ones, candidates int) bool {
		return ones*2 >= candidates
	})
	scrubber := iterate(data, func(ones, candidates int) bool {
		return ones*2 < candidates
	})
	fmt.Println(oxygen * scrubber)
}

func iterate(data [][12]bool, match func(int, int) bool) int {
	potential := make(map[int]bool)
	for i := range data {
		potential[i] = true
	}
	for pos := 0; ; pos++ {
		candidates := len(potential)
		if candidates == 1 {
			var ret int
			for idx := range potential {
				for pos, set := range data[idx] {
					if set {
						ret += 1 << (11 - pos)
					}
				}
			}
			return ret
		}
		oneCount := 0
		for idx := range potential {
			if data[idx][pos] {
				oneCount++
			}
		}
		matched := match(oneCount, candidates)
		for idx := range potential {
			if data[idx][pos] != matched {
				delete(potential, idx)
			}
		}
	}
}
