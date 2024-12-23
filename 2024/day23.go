package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"slices"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day23.input", "Relative file path to use as input.")

type triad [3]string

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	relations := make(map[string][]string)
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, "-")
		a := parts[0]
		b := parts[1]
		relations[a] = append(relations[a], b)
		relations[b] = append(relations[b], a)
	}

	found := make(map[triad]bool)
	for a, v := range relations {
		if len(v) < 2 {
			continue
		}
		if a[0] != 't' {
			continue
		}
		for i := 0; i < len(v)-1; i++ {
			b := v[i]
			for j := i + 1; j < len(v); j++ {
				c := v[j]
				if slices.Contains(relations[b], c) {
					key := []string{a, b, c}
					slices.Sort(key)
					found[triad{key[0], key[1], key[2]}] = true
				}
			}
		}
	}
	fmt.Println(len(found))

	var largest []string
	for a, v := range relations {
		cohort := expand(relations, []string{a}, v)
		if len(cohort) > len(largest) {
			largest = slices.Clone(cohort)
		}
	}
	slices.Sort(largest)
	fmt.Println(strings.Join(largest, ","))
}

func expand(relations map[string][]string, cohort, candidates []string) []string {
	// If we're unable to expand, return just what we have.
	largest := cohort
	proposed := make([]string, len(cohort)+1)
	copy(proposed, cohort)
outer:
	for i, c := range candidates {
		for _, v := range cohort {
			if !slices.Contains(relations[v], c) {
				continue outer
			}
		}
		proposed[len(cohort)] = c
		if res := expand(relations, proposed, candidates[i+1:]); len(res) > len(largest) {
			largest = slices.Clone(res)
		}
	}
	return largest
}
