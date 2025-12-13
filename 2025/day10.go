package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day10.input", "Relative file path to use as input.")

// Assume there won't be more than 64 different lightbulbs. probably fine.
type Model struct {
	Target   uint64
	Len      int
	Buttons  []uint64
	Joltages []int
}

type SequenceA struct {
	Cost     int
	Switches uint64
}

type State [10]uint8

func (m Model) SolveA() int {
	queue := []SequenceA{{0, 0}}
	shortest := make(map[uint64]int)
	// Perform BFS
	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]
		if m.Target == item.Switches {
			return item.Cost
		}
		for _, v := range m.Buttons {
			var next SequenceA
			next.Switches = item.Switches ^ v
			next.Cost = item.Cost + 1
			if lowestCost, ok := shortest[next.Switches]; !ok || next.Cost < lowestCost {
				shortest[next.Switches] = next.Cost
				queue = append(queue, next)
			}
		}
	}
	return -1
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var countA int
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, " ")
		var m Model
		m.Len = len(parts[0]) - 2
		for i, v := range parts[0][1 : len(parts[0])-1] {
			if v == '#' {
				m.Target |= 1 << i
			}
		}
		for _, v := range parts[1 : len(parts)-1] {
			var button uint64
			for _, pos := range strings.Split(v[1:len(v)-1], ",") {
				p, _ := strconv.Atoi(pos)
				button |= 1 << p
			}
			m.Buttons = append(m.Buttons, button)
		}
		joltages := parts[len(parts)-1]
		for _, joltage := range strings.Split(joltages[1:len(joltages)-1], ",") {
			j, _ := strconv.Atoi(joltage)
			m.Joltages = append(m.Joltages, j)
		}
		countA += m.SolveA()
		fmt.Printf("{{")
		for i := range m.Buttons {
			if i != 0 {
				fmt.Printf(",")
			}
			fmt.Printf("a_%d", i)
		}
		fmt.Printf("}} * {")
		for i, v := range m.Buttons {
			if i != 0 {
				fmt.Printf(",")
			}
			fmt.Printf("{")
			for j := 0; j < len(m.Joltages); j++ {
				if j != 0 {
					fmt.Printf(",")
				}
				if v & (1 << j) > 0 {
					fmt.Printf("1")
				} else {
					fmt.Printf("0")
				}
			}
			fmt.Printf("}")
		}
		fmt.Printf("} == {{")
		for i, j := range m.Joltages {
			if i != 0 {
				fmt.Printf(",")
			}
			fmt.Printf("%d", j)
		}
		fmt.Printf("}}\n")
	}
	fmt.Println(countA)
}
