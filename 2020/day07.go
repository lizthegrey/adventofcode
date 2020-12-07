package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day07.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	contains := make(map[string]map[string]int)
	containedBy := make(map[string][]string)

	for _, s := range split {
		parts := strings.Split(s[:len(s)-1], " bags contain ")
		subject := parts[0]
		if parts[1] == "no other bags" {
			continue
		}

		contents := strings.Split(parts[1], ", ")
		contains[subject] = make(map[string]int)
		for _, v := range contents {
			v = strings.TrimSuffix(v, " bags")
			v = strings.TrimSuffix(v, " bag")
			subParts := strings.SplitN(v, " ", 2)
			num, err := strconv.Atoi(subParts[0])
			content := subParts[1]
			if err != nil {
				fmt.Printf("Failed to parse %s\n", s)
				break
			}
			containedBy[content] = append(containedBy[content], subject)
			contains[subject][content] = num
		}
	}

	newlyFound := []string{"shiny gold"}
	found := map[string]bool{
		"shiny gold": true,
	}
	shinyGoldParents := 0
	for {
		nextCycle := make([]string, 0)
		for _, v := range newlyFound {
			validParents := containedBy[v]
			for _, p := range validParents {
				if !found[p] {
					shinyGoldParents++
					found[p] = true
					nextCycle = append(nextCycle, p)
				}
			}
		}

		if len(nextCycle) == 0 {
			break
		}
		newlyFound = nextCycle
	}
	fmt.Println(shinyGoldParents)

	newlyFoundBags := map[string]int{
		"shiny gold": 1,
	}
	allBagsFound := make(map[string]int)
	for {
		nextCycle := make(map[string]int)
		for k, v := range newlyFoundBags {
			children := contains[k]
			for c, count := range children {
				allBagsFound[c] += count * v
				nextCycle[c] += count * v
			}
		}

		if len(nextCycle) == 0 {
			break
		}
		newlyFoundBags = nextCycle
	}
	children := 0
	for _, v := range allBagsFound {
		children += v
	}
	fmt.Println(children)
}
