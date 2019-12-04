package main

import (
	"flag"
	"fmt"
	"strconv"
)

var lower = flag.Int("lower", 123257, "Lower bound.")
var upper = flag.Int("upper", 647015, "Upper bound.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

func main() {
	flag.Parse()
	count := 0
	for i := *lower; i <= *upper; i++ {
		if eligible(i) {
			count++
		}
	}
	fmt.Println(count)
}

func eligible(i int) bool {
	digits := strconv.Itoa(i)
	var prevDigit rune
	var double bool
	var exactDouble bool
	run := 1

	for _, d := range digits {
		if prevDigit == d {
			double = true
			run++
		} else {
			if run == 2 {
				exactDouble = true
			}
			run = 1
		}
		if prevDigit > d {
			return false
		}
		prevDigit = d
	}
	if run == 2 {
		exactDouble = true
	}

	if *partB {
		return exactDouble
	} else {
		return double
	}
}
