package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day21.input", "Relative file path to use as input.")

type monkey struct {
	cached  int
	compute func() int
	pending int
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	// part A
	monkeys := make(map[string]*monkey)
	deps := make(map[string][]string)
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, ": ")
		key := parts[0]

		command := strings.Split(parts[1], " ")
		if len(command) == 1 {
			val, _ := strconv.Atoi(command[0])
			monkeys[key] = &monkey{compute: func() int { return val }}
		} else {
			x := command[0]
			op := command[1]
			y := command[2]

			deps[x] = append(deps[x], key)
			deps[y] = append(deps[y], key)
			monkeys[key] = &monkey{pending: 2, compute: func() int {
				a := monkeys[x].cached
				b := monkeys[y].cached
				switch op {
				case "+":
					return a + b
				case "-":
					return a - b
				case "*":
					return a * b
				case "/":
					return a / b
				}
				return -1
			}}
		}
	}

	for monkeys["root"].pending != -1 {
		for k, m := range monkeys {
			if m.pending != 0 {
				continue
			}
			// This node has newly become unblocked.
			m.pending -= 1
			m.cached = m.compute()
			for _, dep := range deps[k] {
				monkeys[dep].pending -= 1
			}
		}
	}
	fmt.Println(monkeys["root"].cached)
}
