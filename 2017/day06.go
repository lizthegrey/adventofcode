package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

var input = flag.String("input", "14	0	15	12	11	11	3	5	1	6	8	4	9	1	8	4", "The input to use.")

const NumBanks int = 16

type State [NumBanks]int

func main() {
	flag.Parse()

	fields := strings.Fields(*input)
	var allocator State
	seen := make(map[State]int)

	for i, num := range fields {
		n, err := strconv.Atoi(num)
		if err != nil {
			fmt.Printf("Could not parse %s because %v.\n", num, err)
			return
		}
		allocator[i] = n
	}

	i := 1
	for ; seen[allocator] == 0; i++ {
		seen[allocator] = i
		allocator.Reallocate()
	}
	fmt.Printf("%d steps to reach first repeat and %d steps in loop.\n", i-1, i-seen[allocator])
}

func (s *State) Reallocate() {
	highIdx := 0
	highVal := 0
	for i, v := range *s {
		if highVal < v {
			highVal = v
			highIdx = i
		}
	}
	s[highIdx] = 0

	for i := (highIdx + 1) % NumBanks; highVal > 0; i = (i + 1) % NumBanks {
		s[i]++
		highVal--
	}
}
