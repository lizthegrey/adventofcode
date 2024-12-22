package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"slices"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day21.input", "Relative file path to use as input.")

type coord struct {
	r, c int
}

var keypad = map[rune]coord{
	'7': {0, 0},
	'8': {0, 1},
	'9': {0, 2},
	'4': {1, 0},
	'5': {1, 1},
	'6': {1, 2},
	'1': {2, 0},
	'2': {2, 1},
	'3': {2, 2},
	'0': {3, 1},
	'A': {3, 2},
}

var dpad = map[rune]coord{
	'^': {0, 1},
	'A': {0, 2},
	'<': {1, 0},
	'v': {1, 1},
	'>': {1, 2},
}

var keypadHole = coord{3, 0}
var dpadHole = coord{0, 0}

type transition struct {
	before, after rune
	depth         int
}

// combo generates all unique permutations of a, accounting for duplicate characters.
func combo(a []rune, f func([]rune), i int) {
	if i > len(a) {
		f(a)
		return
	}
	combo(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		if a[i] == a[j] {
			continue
		}
		a[i], a[j] = a[j], a[i]
		combo(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}

func (t transition) paths(pad map[rune]coord, hole coord) []string {
	if t.before == t.after {
		return []string{"A"}
	}
	from := pad[t.before]
	to := pad[t.after]

	r := to.r - from.r
	rChar := 'v'
	if r < 0 {
		rChar = '^'
	}

	c := to.c - from.c
	cChar := '>'
	if c < 0 {
		cChar = '<'
	}

	instrs := make([]rune, 0, max(c, -c)+max(r, -r))
	instrs = append(instrs, slices.Repeat([]rune{cChar}, max(c, -c))...)
	instrs = append(instrs, slices.Repeat([]rune{rChar}, max(r, -r))...)

	var ret []string
	combo(instrs, func(in []rune) {
		// check if we pass over a hole.
		loc := from
		for _, v := range in {
			switch v {
			case 'v':
				loc.r++
			case '^':
				loc.r--
			case '>':
				loc.c++
			case '<':
				loc.c--
			}
			if loc == hole {
				return
			}
		}
		seq := make([]rune, len(in)+1)
		copy(seq, in)
		seq[len(in)] = 'A'
		ret = append(ret, string(seq))
	}, 0)
	return ret
}

func dpadManual(memo map[transition]uint64, sequence string) uint64 {
	cur := 'A'
	var ret uint64
	for _, next := range sequence {
		key := transition{cur, next, 0}
		cur = next
		if seq, ok := memo[key]; ok {
			ret += seq
			continue
		}
		// All the paths are equivalent, so just pick one.
		sub := key.paths(dpad, dpadHole)[0]
		memo[key] = uint64(len(sub))
		ret += uint64(len(sub))
	}
	return ret
}

func movements(memo map[transition]uint64, sequence string, depth int, initial bool) uint64 {
	if depth == 0 {
		return dpadManual(memo, sequence)
	}

	cur := 'A'
	var ret uint64
	for _, next := range sequence {
		key := transition{cur, next, depth}
		cur = next
		if seq, ok := memo[key]; ok {
			ret += seq
			continue
		}
		var shortest uint64 = math.MaxUint64
		var candidates []string
		if initial {
			candidates = key.paths(keypad, keypadHole)
		} else {
			candidates = key.paths(dpad, dpadHole)
		}
		for _, candidate := range candidates {
			length := movements(memo, candidate, depth-1, false)
			if length < shortest {
				shortest = length
			}
		}
		memo[key] = shortest
		ret += shortest
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

	var sequences []string
	for _, s := range split[:len(split)-1] {
		sequences = append(sequences, s)
	}

	var sumA, sumB uint64
	memo := make(map[transition]uint64)
	for _, seq := range sequences {
		n, _ := strconv.Atoi(seq[:len(seq)-1])

		stepsA := movements(memo, seq, 2, true)
		complexityA := stepsA * uint64(n)
		sumA += complexityA
	}

	clear(memo)
	for _, seq := range sequences {
		n, _ := strconv.Atoi(seq[:len(seq)-1])

		stepsB := movements(memo, seq, 25, true)
		complexityB := stepsB * uint64(n)
		sumB += complexityB
	}
	fmt.Println(sumA)
	fmt.Println(sumB)
}