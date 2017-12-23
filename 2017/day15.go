package main

import (
	"flag"
	"fmt"
)

var seedA = flag.Int("seedA", 65, "The first seed to use.")
var seedB = flag.Int("seedB", 8921, "The second seed to use.")

var factorA = flag.Int("factorA", 16807, "The first multiplier to use.")
var factorB = flag.Int("factorB", 48271, "The second multiplier to use.")

var modulus = flag.Int("modulus", 2147483647, "The modulus to use.")

var cycles = flag.Int("cycles", 40000000, "The number of cycles to check for matches in.")

var picky = flag.Bool("picky", true, "Whether to only emit values with certain divisibility.")

type State struct {
	A, B uint64
}

func main() {
	flag.Parse()

	initial := State{uint64(*seedA), uint64(*seedB)}
	found := 0
	s := initial
	for i := 0; i < *cycles; i++ {
		s = s.cycle()
		if s.check() {
			found++
		}
	}
	fmt.Printf("Found %d matches.\n", found)
}

func (s State) cycle() State {
	if !*picky {
		return State{
			A: (s.A * uint64(*factorA)) % uint64(*modulus),
			B: (s.B * uint64(*factorB)) % uint64(*modulus),
		}
	} else {
		a := s.A
		b := s.B
		for {
			a = (a * uint64(*factorA)) % uint64(*modulus)
			if a%4 == 0 {
				break
			}
		}
		for {
			b = (b * uint64(*factorB)) % uint64(*modulus)
			if b%8 == 0 {
				break
			}
		}
		return State{a, b}
	}
}

func (s State) check() bool {
	return ((s.A^s.B)&((1<<16)-1) == 0)
}
