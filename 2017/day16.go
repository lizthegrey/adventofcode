package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day16.input", "Relative file path to use as input.")

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

	for _, inst := range instructions {
		switch inst[0] {
		case 's':
			count, err := strconv.Atoi(inst[1:])
			if err != nil {
				fmt.Printf("Failed to parse instruction %s\n", inst)
				return
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
				return
			}
			b, err := strconv.Atoi(operands[1])
			if err != nil {
				fmt.Printf("Failed to parse instruction %s: %s\n", inst)
				return
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
	for _, v := range order {
		fmt.Printf("%c", v)
	}
	fmt.Println()
}
