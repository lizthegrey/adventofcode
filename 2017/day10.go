package main

import (
	"container/ring"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

var input = flag.String("input", "3,4,1,5", "The input to the problem.")
var length = flag.Int("length", 5, "The number of elements in the circular list.")
var rounds = flag.Int("rounds", 1, "The number of rounds to perform.")
var partB = flag.Bool("partB", true, "Whether to perform the ASCII conversion of part B.")

func main() {
	flag.Parse()

	var ls []int
	if !*partB {
		lengths := strings.Split(*input, ",")
		ls = make([]int, len(lengths))
		for i, l := range lengths {
			length, err := strconv.Atoi(l)
			if err != nil {
				fmt.Printf("Failed to parse input.\n")
				return
			}
			ls[i] = length
		}
	} else {
		ls = make([]int, len(*input)+5)
		for i, l := range *input {
			ls[i] = int(l)
		}
		addedLengths := [5]int{17, 31, 73, 47, 23}
		for i := 0; i < len(addedLengths); i++ {
			ls[i+len(*input)] = addedLengths[i]
		}
	}

	loop := ring.New(*length)
	for i := 0; i < *length; i++ {
		Set(loop, i)
		loop = loop.Next()
	}

	pos := loop
	skip := 0
	offset := 0

	for i := 0; i < *rounds; i++ {
		pos = round(pos, ls, &skip, &offset)
	}

	final := pos.Move(-offset)
	if !*partB {
		Debug(final)
		a, b := Get(final), Get(final.Next())
		fmt.Printf("Answer: %d*%d = %d\n", a, b, a*b)
		return
	}

	for i := 0; i < *length; i += 16 {
		intermediate := 0
		for j := 0; j < 16; j++ {
			digit := Get(final)
			final = final.Next()
			intermediate ^= digit
		}
		fmt.Printf("%02x", intermediate)
	}
	fmt.Println()
}

// round performs a round of the "hashing" given the specified constraints.
// Note that this modifies the input ring (danger will robinson)
func round(r *ring.Ring, ls []int, skip, offset *int) *ring.Ring {
	pos := r
	for _, n := range ls {
		if n != 0 {
			leftover := pos.Move(n)
			working := pos.Prev().Unlink(n)
			reversed := Reverse(working)

			if n == *length {
				pos = reversed
			} else {
				pos = leftover.Prev().Link(reversed)
			}
		}
		pos = pos.Move(*skip)
		*offset = (*offset + *skip + n) % *length
		*skip++
	}
	return pos
}

func Set(r *ring.Ring, v int) {
	r.Value = v
}

func Get(r *ring.Ring) int {
	return r.Value.(int)
}

func Debug(r *ring.Ring) {
	pos := r
	for i := 0; i < r.Len(); i++ {
		fmt.Printf("%d,", Get(pos))
		pos = pos.Next()
	}
	fmt.Println()
}

func Reverse(r *ring.Ring) *ring.Ring {
	// [A], B, C, D -> [D], C, B, A
	pos := r.Prev()
	ret := ring.New(r.Len())
	for i := 0; i < r.Len(); i++ {
		Set(ret, Get(pos))
		ret = ret.Next()
		pos = pos.Prev()
	}
	return ret
}
