package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

	devices := make(map[string][]string)
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, ": ")
		var conns []string
		for _, v := range strings.Split(parts[1], " ") {
			conns = append(conns, v)
		}
		devices[parts[0]] = conns
	}
	partA := compute(devices, nil, "you", "out")
	fmt.Println(partA)
	partB := compute(devices, nil, "svr", "dac") *
		compute(devices, nil, "dac", "fft") *
		compute(devices, nil, "fft", "out")
	partB += compute(devices, nil, "svr", "fft") *
		compute(devices, nil, "fft", "dac") *
		compute(devices, nil, "dac", "out")
	fmt.Println(partB)
}

func compute(devices map[string][]string, memo map[string]uint64, state, fin string) uint64 {
	if memo == nil {
		memo = make(map[string]uint64)
	}
	var ret uint64
	if state == fin {
		return 1
	}
	if ret, ok := memo[state]; ok {
		return ret
	}
	for _, v := range devices[state] {
		ret += compute(devices, memo, v, fin)
	}

	memo[state] = ret
	return ret
}
