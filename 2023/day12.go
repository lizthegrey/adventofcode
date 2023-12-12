package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day12.input", "Relative file path to use as input.")

type Tristate int

const (
	Unknown = iota
	Working
	Broken
)

type Row []Tristate
type Memo map[uint32]uint64

// Safe to recursively evaluate; the recursion depth is limited by len(r)
func (r Row) Evaluate(memo Memo, vals []int) uint64 {
	// Evaluate heuristics: we cannot have an impossible combination.
	var knownBroken, knownWorking, unknown int
	for _, v := range r {
		switch v {
		case Working:
			knownWorking++
		case Broken:
			knownBroken++
		case Unknown:
			unknown++
		}
	}

	if len(vals) == 0 {
		if knownBroken == 0 {
			// Everything else has to be set to not broken. Only one way to do that.
			return 1
		}
		// There's at least one known broken, but we don't expect any broken.
		return 0
	}

	// From here on out we have the invariant that len(vals) > 0.
	var minLength int
	var expectedBroken int
	for i, v := range vals {
		expectedBroken += v
		minLength += v
		if i != 0 {
			// To account for space between the runs of brokens.
			minLength++
		}
	}

	if knownBroken > expectedBroken {
		return 0
	}
	if knownBroken+unknown < expectedBroken {
		return 0
	}
	if len(r) < minLength {
		return 0
	}

	brokenSequence := 0
	firstBroken := 0
	for i, v := range r {
		switch v {
		case Working:
			if firstBroken == i {
				firstBroken = i + 1
			} else {
				// This marks the end of the broken sequence.
				if brokenSequence == vals[0] {
					// Valid combination, remove the seen broken sequence
					// both from our search vals and from the prefix of our sequence.
					return r[i:].Evaluate(memo, vals[1:])
				} else {
					// Mismatched, expected != seen
					return 0
				}
			}
		case Broken:
			if brokenSequence >= vals[0] {
				// We aren't allowed to have a sequence longer than brokenSequence.
				return 0
			}
			brokenSequence++
		case Unknown:
			// If we've seen this before, return it.
			key := uint32(len(r)-i)<<16 + uint32(brokenSequence<<8) + uint32(len(vals))
			if cached, ok := memo[key]; ok {
				return cached
			}

			// We're at the first ?. Combine the results for working & broken.
			var sum uint64
			if brokenSequence == 0 {
				sum += r[i+1:].Evaluate(memo, vals)
			} else if brokenSequence == vals[0] {
				sum += r[i+1:].Evaluate(memo, vals[1:])
			}
			if brokenSequence < vals[0] {
				clone := make(Row, len(r)-firstBroken)
				copy(clone, r[firstBroken:])
				clone[i-firstBroken] = Broken
				sum += clone.Evaluate(memo, vals)
			}

			// Then memoise and return.
			memo[key] = sum
			return sum
		}
	}
	// We've gotten to the end of the string. Hopefully we are just correct.
	if len(vals) == 1 && brokenSequence == vals[0] {
		return 1
	}
	// Mismatched count of expected values.
	return 0
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var sumA, sumB uint64
	for _, s := range split[:len(split)-1] {
		var rowA Row
		parts := strings.Split(s, " ")
		for _, v := range parts[0] {
			var t Tristate
			switch v {
			case '.':
				t = Working
			case '#':
				t = Broken
			}
			rowA = append(rowA, t)
		}
		var critA []int
		for _, v := range strings.Split(parts[1], ",") {
			c, err := strconv.Atoi(v)
			if err != nil {
				log.Fatalf("Failure parsing criterion %s: %v", v, err)
			}
			critA = append(critA, c)
		}
		memo := make(Memo)
		sumA += rowA.Evaluate(memo, critA)

		rowB := make(Row, 5*len(rowA)+4)
		critB := make([]int, 5*len(critA))
		for i := 0; i < 5; i++ {
			copy(rowB[i*(len(rowA)+1):], rowA)
			if i != 4 {
				rowB[(i+1)*(len(rowA)+1)-1] = Unknown
			}
			copy(critB[i*len(critA):], critA)
		}
		clear(memo)
		sumB += rowB.Evaluate(memo, critB)
	}
	fmt.Println(sumA)
	fmt.Println(sumB)
}
