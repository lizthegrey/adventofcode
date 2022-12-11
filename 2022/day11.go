package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day11.input", "Relative file path to use as input.")

type monkey struct {
	items    []int
	oper     func(int) int
	next     func(int) int
	activity int
}

type monkeys []monkey

func (ms monkeys) eval(rounds int, reduceWorry func(int) int) {
	for round := 0; round < rounds; round++ {
		for i := range ms {
			for _, worry := range ms[i].items {
				ms[i].activity++
				worry = ms[i].oper(worry)
				worry = reduceWorry(worry)
				next := ms[i].next(worry)
				ms[next].items = append(ms[next].items, worry)
			}
			// All of current monkey's items have been thrown.
			// Clear our items and move to next monkey.
			ms[i].items = nil
		}
	}
}

func (ms monkeys) business() int {
	var top [2]int
	for _, m := range ms {
		if m.activity > top[0] {
			top[1] = top[0]
			top[0] = m.activity
		} else if m.activity > top[1] {
			top[1] = m.activity
		}
	}
	return top[0] * top[1]
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var input monkeys
	lcm := 1
	for r := 0; r < len(split); r += 7 {
		var m monkey
		items := strings.Split(split[r+1][18:], ", ")
		for _, i := range items {
			num, _ := strconv.Atoi(i)
			m.items = append(m.items, num)
		}
		if num, err := strconv.Atoi(split[r+2][25:]); err == nil {
			switch split[r+2][23] {
			case '*':
				m.oper = func(i int) int {
					return i * num
				}
			case '+':
				m.oper = func(i int) int {
					return i + num
				}
			}
		} else {
			switch split[r+2][23] {
			case '*':
				m.oper = func(i int) int {
					return i * i
				}
			case '+':
				m.oper = func(i int) int {
					return i + i
				}
			}
		}
		divisor, _ := strconv.Atoi(strings.Split(split[r+3], " ")[5])
		matched, _ := strconv.Atoi(strings.Split(split[r+4], " ")[9])
		notMatched, _ := strconv.Atoi(strings.Split(split[r+5], " ")[9])
		lcm *= divisor
		m.next = func(i int) int {
			if i%divisor == 0 {
				return matched
			}
			return notMatched
		}
		input = append(input, m)
	}

	// part A
	partA := make(monkeys, len(input))
	copy(partA, input)
	partA.eval(20, func(i int) int {
		return i / 3
	})
	fmt.Println(partA.business())

	// part B
	partB := make(monkeys, len(input))
	copy(partB, input)
	partB.eval(10000, func(i int) int {
		return i % lcm
	})
	fmt.Println(partB.business())
}
