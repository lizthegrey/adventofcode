package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day07.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", true, "Whether to use part B logic.")

type Program struct {
	Name        string
	Weight      int
	Children    []string
	Orphan      bool
	TotalWeight int
	Balanced    bool
}
type Progs map[string]Program

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes)
	lines := strings.Split(contents, "\n")
	programs := make(Progs)

	re := regexp.MustCompile("([a-z]+) \\(([0-9]+)\\)(?: -> ([ a-z,]+))?")

	for _, l := range lines {
		if len(l) == 0 {
			break
		}
		m := re.FindStringSubmatch(l)
		if m == nil {
			fmt.Printf("Failed to parse '%s'\n", l)
			return
		}
		name := m[1]
		weight, err := strconv.Atoi(m[2])
		if err != nil {
			fmt.Printf("Failed to parse weight in '%s'\n", l)
			return
		}
		children := make([]string, 0)
		if len(m[3]) > 0 {
			children = strings.Split(m[3], ", ")
		}
		programs[name] = Program{name, weight, children, true, 0, true}
	}
	for _, v := range programs {
		for _, child := range v.Children {
			c := programs[child]
			c.Orphan = false
			programs[child] = c
		}
	}
	root := ""
	for k, v := range programs {
		if v.Orphan {
			fmt.Printf("Found root node: %s\n", k)
			root = k
		}
	}
	programs.Check(root)
}

func (pgs Progs) Check(n string) {
	p := pgs[n]
	p.TotalWeight = p.Weight
	lastChildWeight := 0
	allChildrenBalanced := true
	for _, c := range p.Children {
		pgs.Check(c)
		childWeight := pgs[c].TotalWeight
		if !pgs[c].Balanced {
			allChildrenBalanced = false
		}
		if lastChildWeight != 0 && lastChildWeight != childWeight {
			p.Balanced = false
		}
		lastChildWeight = childWeight
		p.TotalWeight += childWeight
	}
	if allChildrenBalanced && !p.Balanced {
		// Then the problem is in the direct weights of one of my children.
		// Spit out for a human to analyze.
		fmt.Printf("Found source of imbalance: node %s has imbalanced children: ", n)
		for _, c := range p.Children {
			fmt.Printf("%s -> %d, ", c, pgs[c].TotalWeight)
		}
		fmt.Printf("\n")
	}
	pgs[n] = p
}
