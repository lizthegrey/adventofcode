package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"slices"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day11.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	gates := make(map[string][]string)
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, ": ")
		in := parts[0]
		var out []string
		for _, v := range strings.Split(parts[1], " ") {
			out = append(out, v)
		}
		gates[in] = out
	}
	fmt.Println(compute(gates, []string{"you"}, "out"))
	var countB uint64
	countB += compute(gates, []string{"svr"}, "dac") * compute(gates, []string{"svr", "dac"}, "fft") * compute(gates, []string{"svr", "dac", "fft"}, "out")
	countB += compute(gates, []string{"svr"}, "fft") * compute(gates, []string{"svr", "fft"}, "dac") * compute(gates, []string{"svr", "fft", "dac"}, "out")
	fmt.Println(countB)
}

func compute(gates map[string][]string, visited []string, fin string) uint64 {
	cur := visited[len(visited)-1]
	if cur == fin {
		return 1
	}
	var ret uint64
	visited = slices.Grow(visited, 1)
	for _, v := range gates[cur] {
		if slices.Contains(visited, v) {
			continue
		}
		visited = visited[:len(visited)+1]
		visited[len(visited)-1] = v
		ret += compute(gates, visited, fin)
		visited = visited[:len(visited)-1]
	}
	return ret
}
