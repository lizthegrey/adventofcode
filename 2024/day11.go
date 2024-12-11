package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day11.input", "Relative file path to use as input.")

type memo map[pair]uint64

func (m memo) evolve(v int, gens int) uint64 {
	if gens == 0 {
		return 1
	}

	key := pair{v, gens}
	if ret := m[key]; ret != 0 {
		return ret
	}

	if v == 0 {
		ret := m.evolve(1, gens-1)
		m[key] = ret
		return ret
	}
	var digits int
	for test := v; test > 0; test /= 10 {
		digits++
	}
	if digits%2 == 0 {
		cutoff := 1
		for i := 0; i < digits/2; i++ {
			cutoff *= 10
		}
		left := v / cutoff
		right := v % cutoff
		ret := m.evolve(left, gens-1) + m.evolve(right, gens-1)
		m[key] = ret
		return ret
	}
	ret := m.evolve(v*2024, gens-1)
	m[key] = ret
	return ret
}

type pair struct {
	initial, rounds int
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var stones []int
	for _, s := range strings.Split(split[0], " ") {
		v, _ := strconv.Atoi(s)
		stones = append(stones, v)
	}

	cache := make(memo)
	var sumA, sumB uint64
	for _, v := range stones {
		sumA += cache.evolve(v, 25)
	}
	fmt.Println(sumA)
	for _, v := range stones {
		sumB += cache.evolve(v, 75)
	}
	fmt.Println(sumB)
}
