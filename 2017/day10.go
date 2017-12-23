package main

import (
	"flag"
	"fmt"
	"knot"
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
		ls = knot.Key(*input)
	}

	final := knot.Hash(*length, *rounds, ls)

	if !*partB {
		knot.Debug(final)
		a, b := knot.Get(final), knot.Get(final.Next())
		fmt.Printf("Answer: %d*%d = %d\n", a, b, a*b)
		return
	}

	d := knot.Densify(final)
	for _, v := range d {
		fmt.Printf("%02x", v)
	}
	fmt.Println()
}
