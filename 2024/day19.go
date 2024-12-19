package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day19.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var elements, patterns []string
	for i, s := range split[:len(split)-1] {
		if i == 0 {
			for _, v := range strings.Split(s, ", ") {
				elements = append(elements, v)
			}
		}
		if i < 2 {
			continue
		}
		patterns = append(patterns, s)
	}

	var possible, total uint64
	memo := make(map[string]uint64)
	for _, p := range patterns {
		num := test(memo, elements, p)
		total += num
		if num > 0 {
			possible++
		}
	}
	fmt.Println(possible)
	fmt.Println(total)
}

func test(memo map[string]uint64, elements []string, p string) uint64 {
	if len(p) == 0 {
		return 1
	}
	if v, ok := memo[p]; ok {
		return v
	}
	var total uint64
	for _, e := range elements {
		if len(p) < len(e) {
			continue
		}
		if p[:len(e)] == e {
			total += test(memo, elements, p[len(e):])
		}
	}
	memo[p] = total
	return total
}
