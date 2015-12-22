package main

import (
	"fmt"
)

func main() {
	volumes := []int{7, 10, 11, 18, 18, 21, 22, 24, 26, 32, 36, 40, 40, 42, 43, 44, 46, 47, 49, 50}
	target := 150
	combinations := 1
	for _ = range volumes {
		combinations *= 2
	}

	contCounts := make(map[int]int)

	valid := 0
comb:
	for i := 0; i < combinations; i++ {
		mask := i
		total := 0
		containers := 0
		for v := range volumes {
			if mask%2 == 1 {
				total += volumes[v]
				containers++
			}
			mask >>= 1
			if total > target {
				continue comb
			}
		}
		if total == target {
			valid++
			contCounts[containers]++
		}
	}
	fmt.Println(contCounts)
}
