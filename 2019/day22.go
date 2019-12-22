package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day22.input", "Relative file path to use as input.")

type Deck [10007]int
type Operation func(pos, deckSize int64, iCache map[int64]int64) int64

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}

	var d Deck
	for i := range d {
		d[i] = i
	}

	ops := make([]Operation, 0)

	// Part A
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	for _, s := range split {
		if s == "" {
			continue
		}
		tokens := strings.Split(s, " ")
		var newDeck Deck
		var op Operation
		switch tokens[0] {
		case "cut":
			count, err := strconv.Atoi(tokens[1])
			if err != nil {
				fmt.Printf("Failed to parse cut count %s\n", tokens[1])
			}
			if count > 0 {
				copy(newDeck[len(d)-count:], d[0:count])
				copy(newDeck[0:len(d)-count], d[count:])

				op = func(pos, deckSize int64, _ map[int64]int64) int64 {
					if deckSize-int64(count) <= pos {
						// We were shoved to the end, so we need to reverse.
						// 0123456789 cut 4 ->
						// 4567890123
						return int64(count) - (deckSize - pos)
					}
					// We were shifted left, so we need to shift right.
					return pos + int64(count)
				}
			} else {
				copy(newDeck[0:-count], d[len(d)+count:])
				copy(newDeck[-count:], d[0:len(d)+count])

				op = func(pos, deckSize int64, _ map[int64]int64) int64 {
					// count is negative.
					if pos < -int64(count) {
						// We were shoved to the front, so we need to reverse.
						// 0123456789 cut -4 ->
						// 6789012345
						return deckSize + int64(count) + pos
					}
					// We were shifted right, so we need to shift left.
					return pos + int64(count)
				}
			}
		case "deal":
			switch tokens[1] {
			case "into":
				for i, v := range d {
					newDeck[len(d)-i-1] = v
				}
				op = func(pos, deckSize int64, _ map[int64]int64) int64 {
					// 0123456789
					// 9876543210
					return deckSize - 1 - pos
				}
			case "with":
				inc, err := strconv.Atoi(tokens[3])
				if err != nil {
					fmt.Printf("Failed to parse increment %s\n", tokens[4])
					return
				}
				for i, v := range d {
					newDeck[(inc*i)%len(d)] = v
				}

				op = func(pos, deckSize int64, iCache map[int64]int64) int64 {
					// 0123456789
					// 0741852963
					// As long as count does not divide deckSize evenly (which we were promised),
					// we can reverse division with modulus by multiplying by a different number
					// and take the mod by deckSize as usual.
					// In this example, N is 10, we've sorted by 3
					// To reverse, take current position and multiply 7 and then modulus 10.
					// N=11, sort by 4
					// to reverse, multiply by 3 and then take modulus.
					// 0123456789A
					// 0369147A258
					// the relation we're looking for is A*B === 1 (mod N)
					if _, ok := iCache[int64(inc)]; !ok {
						iCache[int64(inc)] = eea(int64(inc), deckSize)
					}
					var result big.Int
					a := big.NewInt(iCache[int64(inc)])
					b := big.NewInt(pos)
					m := big.NewInt(deckSize)
					result.Mul(a, b)
					result.Mod(&result, m)
					ret := result.Int64()
					return ret
				}
			}
		}
		d = newDeck
		ops = append(ops, op)
	}

	for i, v := range d {
		if v == 2019 {
			fmt.Println(i)
			break
		}
	}

	cCache := make(map[int64]int64)
	for i, v := range d {
		if r := reverse(int64(i), int64(10007), ops, cCache); r != int64(v) {
			fmt.Printf("Failed to match value for position %d: %d != %d\n", i, v, r)
			return
		}
	}

	// Part B:
	// Find the value of the card at position 2020 after 101741582076661 runs
	// on a deck of length 119315717514047.
	// This means we need to reverse the process for one run.
	// Then reverse the reversal.
	// Until we get a repeat and we know we've cycled.

	// Cycle tracking: track the number of reverse cycles it took to wind up
	// in this position.
	positionsSeen := map[int64]int{
		int64(2020): 0,
	}
	currentPos := int64(2020)
	iCache := make(map[int64]int64)

	var o, cycleLength int
	for i := 1; ; i++ {
		if i%1000000 == 0 {
			fmt.Printf("Tested %d reverse cycles.\n", i)
		}
		newPos := reverse(currentPos, int64(119315717514047), ops, iCache)
		if o, ok := positionsSeen[newPos]; ok {
			cycleLength = i - o
			fmt.Printf("Repeat found between cycles %d and %d (cycle length %d).\n", o, i, cycleLength)
			break
		}
		positionsSeen[newPos] = i
		currentPos = newPos
	}
	totalCycles := int64(101741582076661)
	// Once we find our first repeat, then we take total cycles minus the start index.
	// and we take the modulus.
	adjusted := (totalCycles - int64(o)) % int64(cycleLength)
	adjusted += int64(o)
	for k, v := range positionsSeen {
		if int64(v) == adjusted {
			fmt.Println(k)
		}
	}
}

func reverse(pos, deckSize int64, ops []Operation, iCache map[int64]int64) int64 {
	reversedPos := pos
	for i := len(ops) - 1; i >= 0; i-- {
		f := ops[i]
		reversedPos = f(reversedPos, deckSize, iCache)
	}
	return reversedPos
}

func eea(inc, deckSize int64) int64 {
	ret := int64(0)
	newT := int64(1)
	r := deckSize
	newR := inc
	for newR != 0 {
		q := r / newR
		ret, newT = newT, ret-q*newT
		r, newR = newR, r-q*newR
	}

	if r > 1 {
		fmt.Println("Not invertible.")
		return -1
	}
	if ret < 0 {
		ret += deckSize
	}
	return ret
}
