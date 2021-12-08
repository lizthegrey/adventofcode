package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day08.input", "Relative file path to use as input.")

type Entry map[[10]string][4]string

var numberMapping map[string]int = map[string]int{
	"abcefg":  0,
	"cf":      1,
	"acdeg":   2,
	"acdfg":   3,
	"bcdf":    4,
	"abdfg":   5,
	"abdefg":  6,
	"acf":     7,
	"abcdefg": 8,
	"abcdfg":  9,
}

type DecryptionPossibilities map[rune]bool
type DecryptionKey map[rune]DecryptionPossibilities

func (d DecryptionKey) recordResult(scrambled string, golden string) {
	// "ab" => "cf"
	// a => either c or f, eliminate all other possibilities.
	// same for b.
	// and every character except a or b has to go to not c and not f.
	for s, p := range d {
		if strings.IndexRune(scrambled, s) != -1 {
			// Then eliminate every character except ones in golden.
			for c := range p {
				if strings.IndexRune(golden, c) == -1 {
					delete(p, c)
				}
			}
		} else {
			// Eliminate golden characters.
			for _, g := range golden {
				delete(p, g)
			}
		}
	}
}

func (d DecryptionKey) decrypt(in string, missing map[int]bool, numberToSegments map[int]string) *int {
	chars := []rune(in)
	out := make([]rune, len(chars), len(chars))

	foundCount := 0
	var found string
	for m := range missing {
		if len(in) == len(numberToSegments[m]) {
			found = numberToSegments[m]
			foundCount++
		}
	}
	if foundCount == 1 {
		return recognize(found)
	}

	for i, c := range chars {
		substitutions := d[c]
		if len(substitutions) != 1 {
			return nil
		}
		for subs := range substitutions {
			out[i] = subs
		}
	}
	return recognize(string(out))
}

func recognize(in string) *int {
	chars := []rune(in)
	sort.Slice(chars, func(i, j int) bool { return chars[i] < chars[j] })
	if ret, ok := numberMapping[string(chars)]; ok {
		return &ret
	}
	return nil
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	entries := make(Entry)
	partA := 0
	for _, s := range split {
		parts := strings.Split(s, " | ")
		inputs := strings.Split(parts[0], " ")
		outputs := strings.Split(parts[1], " ")
		var left [10]string
		for i, in := range inputs {
			left[i] = in
		}
		var right [4]string
		for i, v := range outputs {
			right[i] = v
			switch len(v) {
			case 2:
				// Digit 1
				fallthrough
			case 4:
				// Digit 4
				fallthrough
			case 3:
				// Digit 7
				fallthrough
			case 7:
				// Digit 8
				partA++
			}
		}
		entries[left] = right
	}
	fmt.Println(partA)

	numberToSegments := make(map[int]string)
	segmentNameToAppearanceCount := make(map[rune]int)
	for k, v := range numberMapping {
		numberToSegments[v] = k
		for _, c := range k {
			segmentNameToAppearanceCount[c]++
		}
	}

	sum := 0
	for left, right := range entries {
		// Represents the possible set of all mappings from wire to wire.
		wireMapping := make(DecryptionKey)
		// Represents which digits have not yet been found.
		missing := make(map[int]bool)
		// Represents which positions still need to be worked.
		workQueue := make(map[int]bool)

		scrambledAppearanceCount := make(map[rune]int)
		for _, scrambled := range left {
			for _, r := range scrambled {
				scrambledAppearanceCount[r]++
			}
		}

		for i := 0; i < 10; i++ {
			missing[i] = true
			workQueue[i] = true
		}
		for c := 'a'; c <= 'g'; c++ {
			possibilities := make(DecryptionPossibilities)
			for p := 'a'; p <= 'g'; p++ {
				if segmentNameToAppearanceCount[p] == scrambledAppearanceCount[c] {
					possibilities[p] = true
				}
			}
			wireMapping[c] = possibilities
		}

		for len(missing) != 0 {
			// Iterate successively until we have recognized all digits.
			for pos := range workQueue {
				scrambled := left[pos]
				if match := wireMapping.decrypt(scrambled, missing, numberToSegments); match != nil {
					wireMapping.recordResult(scrambled, numberToSegments[*match])
					delete(missing, *match)
					delete(workQueue, pos)
					continue
				}
			}
		}

		number := 0
		decimalPlace := 1
		for pos := 0; pos < len(right); pos++ {
			result := wireMapping.decrypt(right[len(right)-1-pos], nil, numberToSegments)
			number += decimalPlace * *result
			decimalPlace *= 10
		}
		sum += number
	}
	fmt.Println(sum)
}
