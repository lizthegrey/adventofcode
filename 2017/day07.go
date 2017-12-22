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
	Name     string
	Weight   int
	Children []string
	Orphan   bool
}

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes)
	lines := strings.Split(contents, "\n")
	programs := make(map[string]Program)

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
		if len(m) == 4 {
			children = strings.Split(m[3], ", ")
		}
		programs[name] = Program{name, weight, children, true}
	}
	for _, v := range programs {
		for _, child := range v.Children {
			c := programs[child]
			c.Orphan = false
			programs[child] = c
		}
	}
	for k, v := range programs {
		if v.Orphan {
			fmt.Printf("Found orphaned node: %s\n", k)
		}
	}
}
