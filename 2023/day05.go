package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day05.input", "Relative file path to use as input.")

type Transform struct {
	Range
	Offset int
}

type Mapping []Transform

type Range struct {
	Min, Len int
}

func (m Mapping) Map(in int) int {
	for _, v := range m {
		if v.Min <= in && in <= v.Min+v.Len {
			return in + v.Offset
		}
	}
	return in
}

func (t Transform) Overlap(in Range) ([2]Range, Range) {
	inMax := in.Min + in.Len
	tMax := t.Min + t.Len

	matchMin := max(in.Min, t.Min)
	matchMax := min(tMax, inMax)

	return [2]Range{
		// Range before match, if exists
		{in.Min, matchMin - in.Min},
		// Range after match, if exists
		{matchMax, inMax - matchMax},
	}, Range{matchMin + t.Offset, matchMax - matchMin}
}

func (m Mapping) MapRange(in Range) []Range {
	var ret []Range
	workset := []Range{in}

outer:
	for len(workset) > 0 {
		item := workset[0]
		workset = workset[1:]
		for _, tf := range m {
			remainder, out := tf.Overlap(item)
			if out.Len <= 0 {
				// Rule did not overlap.
				continue
			}
			ret = append(ret, out)
			for _, r := range remainder {
				if r.Len > 0 {
					workset = append(workset, r)
				}
			}
			// Transform successful, worklist updated.
			continue outer
		}
		// If we got this far, nothing matched. Transfer it over 1 for 1.
		ret = append(ret, item)
	}

	return ret
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var inputA []int
	for _, v := range strings.Split(split[0][7:], " ") {
		seed, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("Failed parsing %s: %v", v, err)
		}
		inputA = append(inputA, seed)
	}
	var inputB []Range
	for i := 0; i < len(inputA); i += 2 {
		inputB = append(inputB, Range{inputA[i], inputA[i+1]})
	}

	var m Mapping
	for _, s := range split[2:len(split)] {
		if s == "" {
			// We've started a new conversion map. Finish processing through.
			var outputA []int
			for _, k := range inputA {
				outputA = append(outputA, m.Map(k))
			}
			inputA = outputA

			var outputB []Range
			for _, k := range inputB {
				outputB = append(outputB, m.MapRange(k)...)
			}
			inputB = outputB

			m = nil
			continue
		}
		if s[len(s)-1] == ':' {
			// This just says what is mapping to what.
			continue
		}

		parts := strings.Split(s, " ")
		if len(parts) != 3 {
			log.Fatalf("Wrong number of parts in line: %s", s)
		}
		dst, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatalf("Failed parsing %s: %v", parts[0], err)
		}
		src, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatalf("Failed parsing %s: %v", parts[1], err)
		}
		length, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Fatalf("Failed parsing %s: %v", parts[2], err)
		}
		m = append(m, Transform{Range{src, length}, dst - src})
	}

	lowestA := -1
	for _, k := range inputA {
		if lowestA == -1 || k < lowestA {
			lowestA = k
		}
	}
	fmt.Println(lowestA)

	lowestB := -1
	for _, k := range inputB {
		if lowestB == -1 || k.Min < lowestB {
			lowestB = k.Min
			// No need to consult length; the min is the lowest in range.
		}
	}
	fmt.Println(lowestB)
}
