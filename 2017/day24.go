package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day24.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", true, "Use length of bridge as primary consideration instead of score.")

type Part struct {
	A, B int
}
type PortMap map[int][]int
type Bridge []int

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes)
	partlist := strings.Split(contents[:len(contents)-1], "\n")

	parts := make([]Part, len(partlist))
	pm := make(PortMap)

	for i, p := range partlist {
		ports := strings.Split(p, "/")
		a, err := strconv.Atoi(ports[0])
		if err != nil {
			fmt.Println("Failed to parse port size.")
		}
		b, err := strconv.Atoi(ports[1])
		if err != nil {
			fmt.Println("Failed to parse port size.")
		}
		parts[i] = Part{a, b}
		if matches, found := pm[a]; found {
			pm[a] = append(matches, i)
		} else {
			pm[a] = []int{i}
		}
		if matches, found := pm[b]; found {
			pm[b] = append(matches, i)
		} else {
			pm[b] = []int{i}
		}
	}

	_, score := Search([]int{}, 0, pm, parts, *partB)
	fmt.Println(score)
}

func Search(incremental Bridge, openPort int, pm PortMap, parts []Part, useLength bool) (int, int) {
	toSearch := pm[openPort]
	longest := 0
	highest := make(map[int]int)
	totalHighest := 0

outer:
	for _, idx := range toSearch {
		for _, used := range incremental {
			if idx == used {
				continue outer
			}
		}
		var newPort int
		if parts[idx].A == openPort {
			newPort = parts[idx].B
		} else {
			newPort = parts[idx].A
		}
		trial := make([]int, len(incremental)+1)
		copy(trial, incremental)
		trial = append(trial, idx)

		length, score := Search(trial, newPort, pm, parts, useLength)

		score += parts[idx].A + parts[idx].B
		length++

		if length > longest {
			longest = length
		}
		if score > highest[length] {
			highest[length] = score
		}
		if score > totalHighest {
			totalHighest = score
		}
	}
	if useLength {
		return longest, highest[longest]
	}
	return longest, totalHighest
}
