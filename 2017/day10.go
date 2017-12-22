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

func main() {
	flag.Parse()

	lengths := strings.Split(*input, ",")
	ls := make([]int, len(lengths))
	for i, l := range lengths {
		length, err := strconv.Atoi(l)
		if err != nil {
			fmt.Printf("Failed to parse input.\n")
			return
		}
		ls[i] = length
	}

	loop := ring.New(*length)
	for i := 0; i < *length; i++ {
		Set(loop, i)
		loop = loop.Next()
	}

	pos := loop
	skip := 0
	offset := 0
	for _, n := range ls {
		Debug(pos)
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
		pos = pos.Move(skip)
		offset = (offset + skip + n) % *length
		skip++
	}
	final := pos.Move(-offset)
	Debug(final)
	a, b := Get(final), Get(final.Next())
	fmt.Printf("Answer: %d*%d = %d\n", a, b, a*b)
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
