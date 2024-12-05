package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"slices"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day05.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	rules := make(map[int][]int)
	var clean int
	var fixed int
	parsingRules := true
	for _, s := range split[:len(split)-1] {
		if len(s) == 0 {
			parsingRules = false
			continue
		}
		if parsingRules {
			parts := strings.Split(s, "|")
			before, _ := strconv.Atoi(parts[0])
			after, _ := strconv.Atoi(parts[1])
			rules[before] = append(rules[before], after)
		} else {
			parts := strings.Split(s, ",")
			var list []int
			for _, v := range parts {
				n, _ := strconv.Atoi(v)
				list = append(list, n)
			}

			cmp := func(a, b int) int {
				for _, v := range rules[b] {
					if v == a {
						return 1
					}
				}
				return -1
			}

			if slices.IsSortedFunc(list, cmp) {
				clean += list[len(list)/2]
			} else {
				slices.SortFunc(list, cmp)
				fixed += list[len(list)/2]
			}
		}
	}
	fmt.Println(clean)
	fmt.Println(fixed)
}
