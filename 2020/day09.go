package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day09.input", "Relative file path to use as input.")
var window = flag.Int("window", 25, "The number of entries in the rolling window.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]
	numbers := make([]int, len(split))
	for i, s := range split {
		n, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("Failed to parse %s\n", s)
			break
		}
		numbers[i] = n
	}

	var weakness int
	for i, v := range numbers {
		if i < *window {
			// preamble.
			continue
		}

		found := false
	middle:
		for n := i - *window; n < i; n++ {
			first := numbers[n]
			if first >= v {
				continue
			}
			for m := n + 1; m < i; m++ {
				second := numbers[m]
				if first+second == v {
					found = true
					break middle
				}
			}
		}
		if !found {
			// This doesn't have something that adds up.
			weakness = v
			break
		}
	}
	fmt.Println(weakness)

	for i, v := range numbers {
		rollingSum := v
		high := v
		low := v
		for n := 1; true; n++ {
			current := numbers[i+n]
			rollingSum += current
			if high < current {
				high = current
			}
			if low > current {
				low = current
			}
			if rollingSum == weakness {
				fmt.Println(high + low)
				return
			}
			if rollingSum > weakness {
				// We busted by going too high; start from the next i.
				break
			}
			// Otherwise keep searching.
		}
	}
}
