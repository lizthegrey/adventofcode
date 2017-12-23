package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day16.input", "Relative file path to use as input.")
var iterations = flag.Int("iterations", 1, "Number of iterations to run.")

type Cache map[[16]byte]int

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes)
	instructions := strings.Split(contents[:len(contents)-1], ",")

	lookup := make(map[byte]int)
	var order [16]byte
	for i := 0; i < 16; i++ {
		order[i] = byte('a') + byte(i)
		lookup[byte('a')+byte(i)] = i
	}

	memo := make(Cache)
	for i := 0; i < *iterations; i++ {
		lookup, order = dance(instructions, lookup, order)
		if prev, found := memo[order]; found {
			loopLen := i - prev
			fmt.Printf("Positions %d and %d repeat (length %d).", prev, i, loopLen)
			break
		} else {
			memo[order] = i
		}
	}

	for _, v := range order {
		fmt.Printf("%c", v)
	}
	fmt.Println()
}

func dance(instructions []string, l map[byte]int, order [16]byte) (map[byte]int, [16]byte) {
	lookup := make(map[byte]int)
	for k, v := range l {
		lookup[k] = v
	}
	for _, inst := range instructions {
		switch inst[0] {
		case 's':
			count, err := strconv.Atoi(inst[1:])
			if err != nil {
				fmt.Printf("Failed to parse instruction %s\n", inst)
				return lookup, order
			}
			var neworder [16]byte
			for i := 0; i < 16; i++ {
				if i < count {
					c := order[16-count+i]
					neworder[i] = c
					lookup[c] = i
				} else {
					c := order[i-count]
					neworder[i] = c
					lookup[c] = i
				}
			}
			order = neworder
		case 'x':
			operands := strings.Split(inst[1:], "/")
			a, err := strconv.Atoi(operands[0])
			if err != nil {
				fmt.Printf("Failed to parse instruction %s: %s\n", inst)
				return lookup, order
			}
			b, err := strconv.Atoi(operands[1])
			if err != nil {
				fmt.Printf("Failed to parse instruction %s: %s\n", inst)
				return lookup, order
			}
			ca := order[a]
			cb := order[b]
			order[a] = cb
			order[b] = ca
			lookup[ca] = b
			lookup[cb] = a
		case 'p':
			ca := inst[1]
			cb := inst[3]
			a := lookup[ca]
			b := lookup[cb]
			order[a] = cb
			order[b] = ca
			lookup[ca] = b
			lookup[cb] = a
		}
	}
	return lookup, order
}
