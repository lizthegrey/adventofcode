package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day07.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")
var workerCount = flag.Int("workers", 5, "The number of parallel workers.")

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	dependencies := make(map[byte][]byte)
	r := bufio.NewReader(f)
	for {
		l, err := r.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]
		temp := strings.Split(l, " ")
		prereq := temp[1][0]
		next := temp[7][0]
		if dependencies[prereq] == nil {
			dependencies[prereq] = []byte{}
		}
		if dependencies[next] == nil {
			dependencies[next] = []byte{prereq}
		} else {
			dependencies[next] = append(dependencies[next], prereq)
		}
	}

	if !*partB {
		sequence := make([]byte, 0)
		workDone := make(map[byte]int)
		for len(sequence) != len(dependencies) {
			candidates := Tick(dependencies, workDone)
			c := candidates[0]
			workDone[c] = 61 + int(c-byte('A'))
			sequence = append(sequence, c)
		}
		fmt.Println(string(sequence))
	} else {
		workDone := make(map[byte]int)
		workers := make([]byte, *workerCount)
		cycles := -1
		for {
			candidates := Tick(dependencies, workDone)
			cycles++

			finished := true
			for k := range dependencies {
				if workDone[k] < 61+int(k-byte('A')) {
					finished = false
				} else {
					for i, inProgress := range workers {
						if inProgress == k {
							workers[i] = 0 // free up worker for new assignment.
						}
					}
				}
			}

			for _, v := range candidates {
				for k, inProgress := range workers {
					if inProgress == 0 {
						workers[k] = v
						workDone[v] = 0
						break
					}
				}
			}
			if finished {
				break
			}
		}
		fmt.Println(cycles)
	}
}

func Tick(dependencies map[byte][]byte, workDone map[byte]int) []byte {
	eligible := make([]byte, 0)

	for k := range workDone {
		workDone[k] += 1
	}

inner:
	for toRun, deps := range dependencies {
		if _, ok := workDone[toRun]; ok {
			continue
		}
		for _, d := range deps {
			if workDone[d] < 61+int(d-byte('A')) {
				continue inner
			}
		}
		eligible = append(eligible, toRun)
	}
	sort.Slice(eligible, func(i, j int) bool { return eligible[i] < eligible[j] })
	return eligible
}
