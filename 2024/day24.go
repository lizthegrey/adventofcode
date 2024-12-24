package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day24.input", "Relative file path to use as input.")

type gate struct {
	a, b, out string
	op        string
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var gateList bool
	var gates []gate
	values := make(map[string]bool)
	for _, s := range split[:len(split)-1] {
		if len(s) == 0 {
			gateList = true
			continue
		}
		if !gateList {
			parts := strings.Split(s, ": ")
			key := parts[0]
			val := parts[1] == "1"
			values[key] = val
		} else {
			parts := strings.Split(s, " ")
			a := parts[0]
			op := parts[1]
			b := parts[2]
			out := parts[4]
			gates = append(gates, gate{a, b, out, op})
		}
	}

	process(gates, values)

	var result uint64
	for i := 0; ; i++ {
		key := fmt.Sprintf("z%02d", i)
		val, ok := values[key]
		if !ok {
			break
		}
		if val {
			result += 1 << i
		}
	}
	fmt.Println(result)

	// Part B is not brute forceable; there are 222 nPr 8 / (2^4 * 4!) possible swaps.
	// Very few people have solved half an hour in, so this one's going to be a toughie.
}

func process(gates []gate, values map[string]bool) {
	for {
		var updated bool
		for _, gate := range gates {
			if _, ok := values[gate.out]; ok {
				continue
			}
			a, ok := values[gate.a]
			if !ok {
				continue
			}
			b, ok := values[gate.b]
			if !ok {
				continue
			}
			updated = true
			switch gate.op {
			case "OR":
				values[gate.out] = a || b
			case "AND":
				values[gate.out] = a && b
			case "XOR":
				values[gate.out] = a != b
			}
		}
		if !updated {
			break
		}
	}
}
