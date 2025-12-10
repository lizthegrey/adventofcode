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

type SequenceB struct {
	Cost    int
	Presses [10]int
}

func (m Model) SolveB() int {
	var empty [10]int
	queue := []SequenceB{{0, empty}}
	shortest := make(map[[10]int]int)

	var final [10]int
	for i, v := range m.Joltages {
		final[i] = v
	}

	// Perform BFS
outer:
	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		if sh, ok := shortest[final]; ok && item.Cost >= sh {
			continue
		}

		finished := true
		for i := range m.Joltages {
			if item.Presses[i] > m.Joltages[i] {
				// Invalid, we exceeded the press count.
				continue outer
			}
			if m.Joltages[i] > item.Presses[i] {
				finished = false
			}
		}
		if finished {
			// Can't argue with perfection.
			continue
		}

		for _, v := range m.Buttons {
			// multiple
			times := -1
			for i := range len(m.Joltages) {
				remaining := m.Joltages[i] - item.Presses[i]
				if v&(1<<i) > 0 {
					if times == -1 || remaining < times {
						times = remaining
					}
				}
			}
			if times > 1 {
				var multiple SequenceB
				multiple.Presses = item.Presses
				for i := range len(m.Joltages) {
					if v&(1<<i) > 0 {
						multiple.Presses[i] += times
					}
				}
				multiple.Cost = item.Cost + times
				if lowestCost, ok := shortest[multiple.Presses]; !ok || multiple.Cost < lowestCost {
					shortest[multiple.Presses] = multiple.Cost
					queue = append(queue, multiple)
				}
			}

			// single
			var single SequenceB
			single.Presses = item.Presses
			for i := range len(m.Joltages) {
				if v&(1<<i) > 0 {
					single.Presses[i]++
				}
			}
			single.Cost = item.Cost + 1
			if lowestCost, ok := shortest[single.Presses]; !ok || single.Cost < lowestCost {
				shortest[single.Presses] = single.Cost
				queue = append(queue, single)
			}
		}
	}
	return shortest[final]
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var countA, countB int
	for _, s := range split[:len(split)-1] {
		fmt.Println(s)
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
		countB += m.SolveB()
	}
	fmt.Println(countA)
	fmt.Println(countB)
}
