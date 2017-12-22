package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day08.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", true, "Whether to use part B logic.")

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes)
	lines := strings.Split(contents, "\n")

	re := regexp.MustCompile("([a-z]+) (inc|dec) ([0-9-]+) if ([a-z]+) (==|<|<=|>|>=|\\!=) ([0-9-]+)")

	regs := make(map[string]int)
	for _, l := range lines {
		if len(l) == 0 {
			break
		}
		m := re.FindStringSubmatch(l)
		if m == nil {
			fmt.Printf("Failed to parse '%s'\n", l)
			return
		}
		r := m[1]
		delta := 1
		if m[2] == "dec" {
			delta = -1
		}
		amount, err := strconv.Atoi(m[3])
		if err != nil {
			fmt.Printf("Failed to parse '%s'\n", l)
			return
		}
		cmp := regs[m[4]]
		operand := m[5]
		value, err := strconv.Atoi(m[6])
		if err != nil {
			fmt.Printf("Failed to parse '%s'\n", l)
			return
		}
		operate := false
		switch operand {
		case "==":
			operate = cmp == value
		case "<":
			operate = cmp < value
		case "<=":
			operate = cmp <= value
		case ">":
			operate = cmp > value
		case ">=":
			operate = cmp >= value
		case "!=":
			operate = cmp != value
		default:
			fmt.Printf("Failed to parse '%s'\n", l)
			return
		}
		if operate {
			regs[r] += delta * amount
		}
	}
	highest := math.MinInt32
	for _, v := range regs {
		if v > highest {
			highest = v
		}
	}
	fmt.Printf("Highest value: %d\n", highest)
}
