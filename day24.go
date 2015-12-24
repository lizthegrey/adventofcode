package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	packages := []int{1, 2, 3, 5, 7, 13, 17, 19, 23, 29, 31, 37, 41, 43, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113}
	numBuckets := 4
	sort.Sort(sort.Reverse(sort.IntSlice(packages)))

	sum := 0
	for i := range packages {
		sum += packages[i]
	}
	max := sum / numBuckets

	foundMatch := false
	bestProduct := uint64(math.MaxUint64)
	for n := 1; !foundMatch; n++ {
		found, product := recursiveTest(packages, max, n)
		if found {
			foundMatch = true
			if product < bestProduct {
				bestProduct = product
			}
		}
	}

	fmt.Println(bestProduct)
}

func recursiveTest(packages []int, maxSum, numPackages int) (bool, uint64) {
	found := false
	bestProduct := uint64(math.MaxUint64)
	for i := range packages {
		max := maxSum - packages[i]
		if max < 0 {
			continue
		} else if max == 0 {
			return true, uint64(packages[i])
		} else if numPackages > 1 && i+i != len(packages) {
			subFound, subProduct := recursiveTest(packages[i+1:], max, numPackages-1)
			if subFound {
				found = true
				product := subProduct * uint64(packages[i])
				if product < bestProduct {
					bestProduct = product
				}
			}
		}
	}
	return found, bestProduct
}
